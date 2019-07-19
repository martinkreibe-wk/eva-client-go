// Copyright 2018-2019 Workiva Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"github.com/Workiva/eva-client-go/edn"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strconv"
	"time"

	"github.com/Workiva/eva-client-go/eva"

	_ "github.com/Workiva/eva-client-go/eva/http"
)

type httpTester struct {
	source eva.Source
	label  string
}

func newTester(server, port, category, label, tenant string) *httpTester {

	config, err := eva.NewConfiguration(fmt.Sprintf(`{
			"source": {
				"retries": "10@5000",
				"type":    "http",
				"server":  "%s:%s",
				"mime":    "application/vnd.eva+edn"
			},

			"category": "%s"
		}`, server, port, category))
	Ω(err).Should(BeNil())
	Ω(config).ShouldNot(BeNil())

	ten, err := eva.NewTenant(tenant)
	Ω(err).Should(BeNil())
	Ω(ten).ShouldNot(BeNil())

	source, err := eva.NewSource(config, ten)
	Ω(err).Should(BeNil())
	Ω(source).ShouldNot(BeNil())

	return &httpTester{
		source: source,
		label:  label,
	}
}

func (tester *httpTester) transact(format string, args ...interface{}) (r *transactResult) {
	trx := fmt.Sprintf(format, args...)

	var conn eva.ConnectionChannel
	conn, err := tester.source.Connection(tester.label)
	Ω(err).Should(BeNil())
	Ω(conn).ShouldNot(BeNil())

	result, err := conn.Transact(trx)
	Ω(err).Should(BeNil())
	Ω(result).ShouldNot(BeNil())

	v, h := result.String()
	Ω(v).ShouldNot(HaveLen(0))
	Ω(h).Should(BeTrue())

	var elem edn.Element
	if elem, err = edn.ParseString(v); err == nil {

		r = &transactResult{
			rawResult:          v,
			partitionedTempIds: make(map[string]map[int64]int64),
			rawTempIds:         make(map[int64]int64),
		}

		if collElem, is := elem.(edn.CollectionElement); is {
			err = collElem.IterateChildren(func(key edn.Element, value edn.Element) (e error) {

				if key.ElementType() == edn.KeywordType {
					switch keyElem := key.(edn.SymbolElement); keyElem.Name() {
					case "tempids":
						if keyElem.Prefix() == "eva.client.service" {
							e = value.(edn.CollectionElement).IterateChildren(func(tempIdKey edn.Element, tempIdValue edn.Element) (e2 error) {
								if tempIdKey.Tag() == "db/id" && tempIdKey.ElementType() == edn.VectorType && tempIdValue.ElementType() == edn.IntegerType {

									keyColl := tempIdKey.(edn.CollectionElement)

									var part edn.Element
									if part, e2 = keyColl.Get(0); e2 == nil {
										p := part.String()
										var id edn.Element
										if id, e2 = keyColl.Get(1); e2 == nil {

											if _, has := r.partitionedTempIds[p]; !has {
												r.partitionedTempIds[p] = make(map[int64]int64)
											}
											r.partitionedTempIds[p][id.Value().(int64)] = tempIdValue.Value().(int64)
										}
									}
								} else {
									e2 = fmt.Errorf("unexpected `eva.client.service/tempids` key or value: `%+v` : `%+v`", tempIdKey, tempIdValue)
								}
								return e2
							})
						} else {
							if value.ElementType() == edn.MapType {
								e = value.(edn.CollectionElement).IterateChildren(func(tempIdKey edn.Element, tempIdValue edn.Element) (e2 error) {
									if tempIdKey.ElementType() == edn.IntegerType && tempIdValue.ElementType() == edn.IntegerType {
										r.rawTempIds[tempIdKey.Value().(int64)] = tempIdValue.Value().(int64)
									} else {
										e2 = fmt.Errorf("unexpected `tempids` key or value: `%+v` : `%+v`", tempIdKey, tempIdValue)
									}
									return e2
								})
							} else {
								e = fmt.Errorf("unexpected `tempids` to be a map")
							}
						}
					case "db-before":
						if value.Tag() == "eva.client.service/snapshot-ref" && value.ElementType() == edn.MapType {
							e = value.(edn.CollectionElement).IterateChildren(func(k edn.Element, v edn.Element) (e2 error) {

								if key.ElementType() == edn.KeywordType {
									switch ke := k.(edn.SymbolElement); ke.Name() {
									case "label":
										r.beforeLabel = v.String()
									case "as-of":
										r.beforeT = v.Value().(int64)
									default:
										e2 = fmt.Errorf("unexpected keys in the `db-before` map, got: %s", k.String())
									}
								} else {
									e2 = fmt.Errorf("expected only keyword keys in the `db-before` map, got: %s", k.String())
								}
								return e2
							})
						} else {
							e = fmt.Errorf("unexpected `db-before` to be a map with the `eva.client.service/snapshot-ref` tag")
						}
					case "db-after":
						if value.Tag() == "eva.client.service/snapshot-ref" && value.ElementType() == edn.MapType {
							e = value.(edn.CollectionElement).IterateChildren(func(k edn.Element, v edn.Element) (e2 error) {

								if key.ElementType() == edn.KeywordType {
									switch ke := k.(edn.SymbolElement); ke.Name() {
									case "label":
										r.afterLabel = v.String()
									case "as-of":
										r.afterT = v.Value().(int64)
									default:
										e2 = fmt.Errorf("unexpected keys in the `db-after` map, got: %s", k.String())
									}
								} else {
									e2 = fmt.Errorf("expected only keyword keys in the `db-after` map, got: %s", k.String())
								}
								return e2
							})
						} else {
							e = fmt.Errorf("unexpected `db-before` to be a map with the `eva.client.service/snapshot-ref` tag")
						}
					case "tx-data":
						if value.ElementType() == edn.ListType {
							e = value.(edn.CollectionElement).IterateChildren(func(_ edn.Element, v edn.Element) (e2 error) {

								if v.Tag() == "datom" && v.ElementType() == edn.VectorType {

									vAsColl := v.(edn.CollectionElement)

									var datomE edn.Element
									var datomA edn.Element
									var datomV edn.Element
									var datomT edn.Element
									var datomR edn.Element

									if datomE, e2 = vAsColl.Get(0); e2 == nil {
										if datomA, e2 = vAsColl.Get(1); e2 == nil {
											if datomV, e2 = vAsColl.Get(2); e2 == nil {
												if datomT, e2 = vAsColl.Get(3); e2 == nil {
													if datomR, e2 = vAsColl.Get(4); e2 == nil {
														r.datoms = append(r.datoms, &datom{
															e: datomE.Value().(int64),
															a: datomA.Value().(int64),
															v: datomV,
															t: datomT.Value().(int64),
															r: datomR.Value().(bool),
														})
													}
												}
											}
										}
									}

								} else {
									e2 = fmt.Errorf("expected vectors with `datom` tags in `tx-data`, Got: %s", v.String())
								}

								return e2
							})
						} else {
							e = fmt.Errorf("unexpected `tx-data` to be a list")
						}
					default:
						e = fmt.Errorf("unexpected key value: %s", key.String())
					}
				} else {
					e = fmt.Errorf("expected only keyword keys in the map, got: %s", key.String())
				}

				return e
			})
		} else {
			err = fmt.Errorf("expected a collection, got: %s", v)
		}
	}

	Ω(err).Should(BeNil())
	Ω(r).ShouldNot(BeNil())
	return r
}

func (tester *httpTester) query(query string, items ...interface{}) string {
	result, err := tester.source.Query(query, items...)
	Ω(err).Should(BeNil())
	Ω(result).ShouldNot(BeNil())

	v, h := result.String()
	Ω(v).ShouldNot(HaveLen(0))
	Ω(h).Should(BeTrue())
	return v
}

func (tester *httpTester) pull(pattern string, ids interface{}, items ...interface{}) string {
	snap, err := tester.source.LatestSnapshot(tester.label)
	Ω(err).Should(BeNil())
	Ω(snap).ShouldNot(BeNil())

	var result eva.Result
	result, err = snap.Pull(pattern, ids, items...)
	Ω(err).Should(BeNil())
	Ω(result).ShouldNot(BeNil())

	v, h := result.String()
	Ω(v).ShouldNot(HaveLen(0))
	Ω(h).Should(BeTrue())
	return v
}

func (tester *httpTester) pullAt(t interface{}, pattern string, ids interface{}, items ...interface{}) string {
	snap, err := tester.source.AsOfSnapshot(tester.label, t)
	Ω(err).Should(BeNil())
	Ω(snap).ShouldNot(BeNil())

	var result eva.Result
	result, err = snap.Pull(pattern, ids, items...)
	Ω(err).Should(BeNil())
	Ω(result).ShouldNot(BeNil())

	v, h := result.String()
	Ω(v).ShouldNot(HaveLen(0))
	Ω(h).Should(BeTrue())
	return v
}

type datom struct {
	e int64
	a int64
	v edn.Element
	t int64
	r bool
}

type transactResult struct {
	rawResult          string
	partitionedTempIds map[string]map[int64]int64
	rawTempIds         map[int64]int64
	beforeLabel        string
	beforeT            int64
	afterLabel         string
	afterT             int64
	datoms             []*datom
}

func (d *transactResult) mayHaveValue(value string) {

	has := false
	var values []interface{}
	for _, datom := range d.datoms {
		if datom.v.ElementType() == edn.StringType {
			if value == datom.v.Value().(string) {
				has = true
			}
		}

		values = append(values, datom.v.Value())
	}

	if !has {
		fmt.Printf(fmt.Sprintf("\n[WARN] Didn't have value in datoms: `%s` instead has: %v\n", value, values))
	}
}

func (d *transactResult) mustHaveValue(value string) {

	has := false
	var values []interface{}
	for _, datom := range d.datoms {
		if datom.v.ElementType() == edn.StringType {
			if value == datom.v.Value().(string) {
				has = true
			}
		}

		values = append(values, datom.v.Value())
	}

	if !has {
		Fail(fmt.Sprintf("Didn't have value in datoms: `%s` instead has: %v", value, values))
	}
}

func (d *transactResult) resultTemp(partition string, id int64) (out int64) {
	if coll, has := d.partitionedTempIds[partition]; has {
		if out, has = coll[id]; !has {
			Fail(fmt.Sprintf("Expected to have a temp result is [%d] in partition: `%s` but didn't", id, partition))
		}
	} else {
		Fail(fmt.Sprintf("Expected to have a temp result partition: `%s` but didn't", partition))
	}

	return out
}

func (d *transactResult) dbAfterT() int64 {
	return d.afterT
}

var _ = Describe("General integration tests", func() {

	if os.Getenv("EVA_TEST_INTEGRATION") == "true" {

		cat := generateCategoryName()
		label := "label"
		host := "localhost"
		port := "8080"
		tenant := "tenant"

		Context("transact the scenario", func() {
			It("should execute with the vector reference representation", func() {

				var transactResult *transactResult
				var err error
				var result string

				fmt.Printf("\n\tTest Category: `%s`\n", cat)
				t := newTester(host, port, cat, label, tenant)

				transactResult = t.transact(BookSchema)
				transactResult.mayHaveValue("Title of a book")
				transactResult.mayHaveValue("Date book was published")
				transactResult.mayHaveValue("Author of a book")
				transactResult.mayHaveValue("Name of author")

				transactResult = t.transact(AddFirstBook)
				transactResult.mustHaveValue("First Book")
				transactResult.mustHaveValue("James Madison")

				elementId := transactResult.resultTemp(":db.part/user", -1)
				firstT := transactResult.dbAfterT()

				transactResult = t.transact(UpdateFirstBookAuthorFormat, elementId)
				transactResult.mustHaveValue("Gilgamesh")

				var snap eva.Reference
				snap, err = eva.NewSnapshotAsOfReference(label, firstT)
				Ω(err).Should(BeNil())
				Ω(snap).ShouldNot(BeNil())

				result = t.query(QueryForAuthor, snap, edn.NewStringElement("First Book"))
				Ω(result).Should(BeEquivalentTo("\"James Madison\""))

				snap, err = eva.NewSnapshotReference(label)
				Ω(err).Should(BeNil())
				Ω(snap).ShouldNot(BeNil())

				result = t.query(QueryForAuthor, snap, edn.NewStringElement("First Book"))
				Ω(result).Should(BeEquivalentTo("\"Gilgamesh\""))

				transactResult = t.transact(AddManyBooks)

				transactResult.mustHaveValue("Martin Kleppman")
				transactResult.mustHaveValue("Designing Data-Intensive Applications")
				transactResult.mustHaveValue("Aurelien Geron")
				transactResult.mustHaveValue("Hands-On Machine Learning")
				transactResult.mustHaveValue("Wil van der Aalst")
				transactResult.mustHaveValue("Process Mining: Data Science in Action")
				transactResult.mustHaveValue("Modeling Business Processes: A Petri-Net Oriented Approach")
				transactResult.mustHaveValue("Designing Data-Intensive Applications")
				transactResult.mustHaveValue("Edward Tufte")
				transactResult.mustHaveValue("The Visual Display of Quantitative Information")
				transactResult.mustHaveValue("Envisioning Information")
				transactResult.mustHaveValue("Designing Data-Intensive Applications")
				transactResult.mustHaveValue("Ramez Elmasri")
				transactResult.mustHaveValue("Operating Systems: A Spiral Approach")
				transactResult.mustHaveValue("Fundamentals of Database Systems")
				transactResult.mustHaveValue("Steve McConnell")
				transactResult.mustHaveValue("Code Complete: A Practical Handbook of Software Construction")
				transactResult.mustHaveValue("Software Estimation: Demystifying the Black Art")
				transactResult.mustHaveValue("Rapid Development: Taming Wild Software Schedules")
				transactResult.mustHaveValue("Software Project Survival Guide")
				transactResult.mustHaveValue("After the Gold Rush: Creating a True Profession of Software Engineering")
				transactResult.mustHaveValue("Don Miguel Ruiz")
				transactResult.mustHaveValue("Charles Petzold")
				transactResult.mustHaveValue("Code: The Hidden Language of Computer Hardware and Software")
				transactResult.mustHaveValue("Anil Maheshwari")
				transactResult.mustHaveValue("Data Analytics Made Accessible")
				transactResult.mustHaveValue("Jeremy Anderson")
				transactResult.mustHaveValue("Professional Clojure")

				result = t.query(QueryForTitlesFrom2017, snap)
				Ω(result).Should(BeEquivalentTo("[[\"Designing Data-Intensive Applications\"] [\"Hands-On Machine Learning\"]]"))

				result = t.query(QueryForTitleDesigningDataIntensiveApplications, snap)
				Ω(result).Should(BeEquivalentTo("[[\"Martin Kleppman\"]]"))

				result = t.query(QueryForTitleFromSteveMcConnell, snap)
				Ω(result).Should(BeEquivalentTo(`[["Code Complete: A Practical Handbook of Software Construction"] ["Software Estimation: Demystifying the Black Art"] ["Software Project Survival Guide"] ["After the Gold Rush: Creating a True Profession of Software Engineering"] ["Rapid Development: Taming Wild Software Schedules"]]`))

				result = t.query(QueryForBookFromTitle, snap, "\"First Book\"")

				var id int
				id, err = strconv.Atoi(result)
				Ω(err).Should(BeNil())

				result = t.pull(PullAll, id)
				Ω(result).Should(ContainSubstring(":author/name \"Gilgamesh\""))
				Ω(result).Should(ContainSubstring(":book/title \"First Book\""))

				result = t.pullAt(firstT, PullAll, id)
				Ω(result).Should(ContainSubstring(":author/name \"James Madison\""))
				Ω(result).Should(ContainSubstring(":book/title \"First Book\""))

				result = t.pull(PullAll, QueryAllBooksFrom2017, snap)
				Ω(result).Should(ContainSubstring(":book/title \"Designing Data-Intensive Applications\""))
				Ω(result).Should(ContainSubstring(":book/title \"Hands-On Machine Learning\""))

				result = t.query(QueryForAllTitles, snap)
				Ω(result).Should(BeEquivalentTo(`[["Code Complete: A Practical Handbook of Software Construction"] ["The Four Agreements: A Practical Guide to Personal Freedom"] ["Designing Data-Intensive Applications"] ["Hands-On Machine Learning"] ["Software Estimation: Demystifying the Black Art"] ["Code: The Hidden Language of Computer Hardware and Software"] ["Modeling Business Processes: A Petri-Net Oriented Approach"] ["Process Mining: Data Science in Action"] ["Operating Systems: A Spiral Approach"] ["Envisioning Information"] ["Data Analytics Made Accessible"] ["Rapid Development: Taming Wild Software Schedules"] ["Fundamentals of Database Systems"] ["The Visual Display of Quantitative Information"] ["Professional Clojure"] ["Software Project Survival Guide"] ["First Book"] ["After the Gold Rush: Creating a True Profession of Software Engineering"]]`))

				result = t.query(QueryForAddTransactionForProcessMining, snap)

				today := fmt.Sprintf("#inst \"%s", time.Now().UTC().Format("2006-01-02T15:04")) // this should have happened about now?
				Ω(result).Should(HavePrefix(today))

				result = t.query(QueryForBooksBefore2005, snap)
				Ω(result).Should(BeEquivalentTo(`[["After the Gold Rush: Creating a True Profession of Software Engineering" 1999] ["Code Complete: A Practical Handbook of Software Construction" 2004] ["The Visual Display of Quantitative Information" 2001] ["Code: The Hidden Language of Computer Hardware and Software" 2000] ["Rapid Development: Taming Wild Software Schedules" 1996] ["Software Project Survival Guide" 1997] ["Envisioning Information" 1990]]`))

				result = t.query(QueryForBooksBeforeSurvivalGuide, snap)
				Ω(result).Should(BeEquivalentTo(`[["Rapid Development: Taming Wild Software Schedules" 1996] ["Envisioning Information" 1990]]`))

				result = t.query(QueryForOldestBookYear, snap)
				Ω(result).Should(BeEquivalentTo("1990"))

				var rules edn.Element
				rules, err = edn.ParseString(DefineBookAuthorRules)
				Ω(err).Should(BeNil())

				result = t.query(QueryForBookAuthorOfPetriNetOrientedApproach, snap, rules)

				Ω(result).Should(BeEquivalentTo("\"Wil van der Aalst\""))

				result = t.query(QueryForBookAuthorOfPetriNetOrientedApproach, snap, DefineBookAuthorRules)

				Ω(result).Should(BeEquivalentTo("\"Wil van der Aalst\""))

				transactResult = t.transact(AddJsonJasonBook)
				transactResult.mustHaveValue("Json Book")
				transactResult.mustHaveValue("Jason")

				result = t.query(QueryForBookFromTitle, snap, "\"Json Book\"")

				id, err = strconv.Atoi(result)
				Ω(err).Should(BeNil())

				result = t.pull(PullAll, id)
				Ω(result).Should(ContainSubstring(":book/title \"Json Book\""))

			})
		})
	}
})
