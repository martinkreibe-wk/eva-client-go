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

const (
	BookSchema = `[ ; [1] - BookSchema
	{
		:db/id #db/id [:db.part/user]
		:db/ident :book/title
		:db/doc "Title of a book"
		:db/valueType :db.type/string
		:db/cardinality :db.cardinality/one
		:db.install/_attribute :db.part/db
	}
	{
		:db/id #db/id [:db.part/user]
		:db/ident :book/year_published
		:db/doc "Date book was published"
		:db/valueType :db.type/long
		:db/cardinality :db.cardinality/one
		:db.install/_attribute :db.part/db
	}
	{
		:db/id #db/id [:db.part/user]
		:db/ident :book/author
		:db/doc "Author of a book"
		:db/valueType :db.type/ref
		:db/cardinality :db.cardinality/one
		:db.install/_attribute :db.part/db
	}
	{
		:db/id #db/id [:db.part/user]
		:db/ident :author/name
		:db/doc "Name of author"
		:db/valueType :db.type/string
		:db/cardinality :db.cardinality/one
		:db.install/_attribute :db.part/db
	}
]`

	AddFirstBook = `[ ; [2] - AddFirstBook
	[
		:db/add #db/id [:db.part/user -1]
		:book/title "First Book"
	]
	[
		:db/add #db/id [:db.part/user -1]
		:author/name "James Madison"
	]
]`

	UpdateFirstBookAuthorFormat = `[ ; [3] - UpdateFirstBookAuthorFormat
	[
		:db/add %d
		:author/name "Gilgamesh"
	]
]`

	QueryForAuthor = `[ ; [4,5] - QueryForAuthor
	:find ?a .
	:in $ ?t
	:where
		[?b :book/title ?t]
		[?b :author/name ?a]
]`

	AddManyBooks = `[ ; [6] - AddManyBooks
	{:db/id #db/id [:db.part/user -1] :author/name "Martin Kleppman"}
	{:db/id #db/id [:db.part/user -2] :book/title "Designing Data-Intensive Applications" :book/year_published 2017 :book/author #db/id[:db.part/user -1]}
	{:db/id #db/id [:db.part/user -3] :author/name "Aurelien Geron"}
	{:db/id #db/id [:db.part/user -4] :book/title "Hands-On Machine Learning" :book/year_published 2017 :book/author #db/id[ :db.part/user -3]}
	{:db/id #db/id [:db.part/user -5] :author/name "Wil van der Aalst"}
	{:db/id #db/id [:db.part/user -6] :book/title "Process Mining: Data Science in Action" :book/year_published 2016 :book/author #db/id[ :db.part/user -5]}
	{:db/id #db/id [:db.part/user -7] :book/title "Modeling Business Processes: A Petri-Net Oriented Approach" :book/year_published 2011 :book/author #db/id[ :db.part/user -5]}
	{:db/id #db/id [:db.part/user -8] :author/name "Edward Tufte"}
	{:db/id #db/id [:db.part/user -9] :book/title "The Visual Display of Quantitative Information" :book/year_published 2001 :book/author #db/id[ :db.part/user -8]}
	{:db/id #db/id [:db.part/user -10] :book/title "Envisioning Information" :book/year_published 1990 :book/author #db/id[ :db.part/user -8]}
	{:db/id #db/id [:db.part/user -11] :author/name "Ramez Elmasri"}
	{:db/id #db/id [:db.part/user -12] :book/title "Operating Systems: A Spiral Approach" :book/year_published 2009 :book/author #db/id[ :db.part/user -11]}
	{:db/id #db/id [:db.part/user -13] :book/title "Fundamentals of Database Systems" :book/year_published 2006 :book/author #db/id[ :db.part/user -11]}
	{:db/id #db/id [:db.part/user -14] :author/name "Steve McConnell"}
	{:db/id #db/id [:db.part/user -15] :book/title "Code Complete: A Practical Handbook of Software Construction" :book/year_published 2004 :book/author #db/id[:db.part/user -14]}
	{:db/id #db/id [:db.part/user -16] :book/title "Software Estimation: Demystifying the Black Art" :book/year_published 2006 :book/author #db/id[ :db.part/user -14]}
	{:db/id #db/id [:db.part/user -17] :book/title "Rapid Development: Taming Wild Software Schedules" :book/year_published 1996 :book/author #db/id[:db.part/user -14]}
	{:db/id #db/id [:db.part/user -18] :book/title "Software Project Survival Guide" :book/year_published 1997 :book/author #db/id[ :db.part/user -14]}
	{:db/id #db/id [:db.part/user -19] :book/title "After the Gold Rush: Creating a True Profession of Software Engineering" :book/year_published 1999 :book/author #db/id[ :db.part/user -14]}
	{:db/id #db/id [:db.part/user -20] :author/name "Don Miguel Ruiz"}
	{:db/id #db/id [:db.part/user -21] :book/title "The Four Agreements: A Practical Guide to Personal Freedom" :book/year_published 2011 :book/author #db/id[ :db.part/user -20]}
	{:db/id #db/id [:db.part/user -22] :author/name "Charles Petzold"}
	{:db/id #db/id [:db.part/user -23] :book/title "Code: The Hidden Language of Computer Hardware and Software" :book/year_published 2000 :book/author #db/id[ :db.part/user -22]}
	{:db/id #db/id [:db.part/user -24] :author/name "Anil Maheshwari"}
	{:db/id #db/id [:db.part/user -25] :book/title "Data Analytics Made Accessible" :book/year_published 2014 :book/author #db/id[ :db.part/user -24]}
	{:db/id #db/id [:db.part/user -26] :author/name "Jeremy Anderson"}
	{:db/id #db/id [:db.part/user -27] :book/title "Professional Clojure" :book/year_published 2016 :book/author #db/id[:db.part/user -26]}
]`

	QueryForTitlesFrom2017 = `[ ; [7] - QueryForTitlesFrom2017
	:find ?title
	:where
		[?b :book/year_published 2017]
		[?b :book/title ?title]
]`

	QueryForTitleDesigningDataIntensiveApplications = `[ ; [8] - QueryForTitle_DesigningDataIntensiveApplications
	:find ?name
	:where
		[?b :book/title "Designing Data-Intensive Applications"]
		[?b :book/author ?a]
		[?a :author/name ?name]
]`

	QueryForTitleFromSteveMcConnell = `[ ; [9] - QueryForTitleFromSteveMcConnell
	:find ?books
	:where
		[?b :book/title ?books]
		[?b :book/author ?a]
		[?a :author/name "Steve McConnell"]
]`

	QueryForBookFromTitle = `[ ; [10,24] - QueryForBookFromTitle
	:find ?b .
	:in $ ?t
	:where [?b :book/title ?t]
]`

	PullAll = "[*] ; [11,12,25] - PullAll"

	QueryAllBooksFrom2017 = "[:find ?b :where [?b :book/year_published 2017] [?b :book/title _]] ; [13] - QueryAllBooksFrom2017"

	QueryForAllTitles = `[ ; [14] - QueryForAllTitles
	:find ?name
	:where
		[_ :book/title ?name]
]`

	QueryForAddTransactionForProcessMining = `[ ; [15] - QueryForAddTransactionForProcessMining
	:find ?timestamp .
	:where
		[_ :book/title "Process Mining: Data Science in Action" ?tx]
		[?tx :db/txInstant ?timestamp]
]`

	QueryForBooksBefore2005 = `[ ; [16] - QueryForBooksBefore2005
	:find ?book ?year
	:where
		[?b :book/title ?book]
		[?b :book/year_published ?year]
		[(< ?year 2005)]
]`

	QueryForBooksBeforeSurvivalGuide = `[ ; [17] - QueryForBooksBeforeSurvivalGuide
	:find ?book ?y1
	:where
		[?b1 :book/title ?book]
		[?b1 :book/year_published ?y1]
		[?b2 :book/title "Software Project Survival Guide"]
		[?b2 :book/year_published ?y2]
		[(< ?y1 ?y2)]
]`

	QueryForOldestBookYear = `[ ; [18] - QueryForOldestBookYear
	:find (min ?year) .
	:where
		[_ :book/year_published ?year]			
]`

	DefineBookAuthorRules = `[ ; [19,21] - DefineBookAuthorRules
	[(book-author ?book ?name)
		[?b :book/title ?book]
		[?b :book/author ?a]
		[?a :author/name ?name]
	]
]`

	QueryForBookAuthorOfPetriNetOrientedApproach = `[ ; [20,22] - QueryForBookAuthorOfPetriNetOrientedApproach
   	:find ?name .
	:in $ %
	:where
		(book-author "Modeling Business Processes: A Petri-Net Oriented Approach" ?name)
]`

	AddJsonJasonBook = `[ ; [23] - AddJsonJasonBook
	[:db/add #db/id [:db.part/user] :book/title "Json Book"]
	[:db/add #db/id [:db.part/user] :author/name "Jason"]
]`
)
