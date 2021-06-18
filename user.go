package main

import (
	"os"

	"github.com/CheatCoder/geddit"
)

func getUserList(npage int) []*geddit.Submission {
	var list []*geddit.Submission
	var err error
	if npage > page {
		list, err = o.MySavedLinks(geddit.ListingOptions{
			After: after,
			Limit: 5,
			Time:  "New",
		})
		if err != nil {
			os.Exit(-50)
		}
	} else if npage < page {
		list, err = o.MySavedLinks(geddit.ListingOptions{
			Before: before,
			Limit:  5,
			Time:   "New",
		})
		if err != nil {
			os.Exit(-50)
		}
	} else {
		list, err = o.MySavedLinks(geddit.ListingOptions{
			Limit: 5,
			Time:  "New",
		})
		if err != nil {
			os.Exit(-50)
		}
	}

	return list
}
