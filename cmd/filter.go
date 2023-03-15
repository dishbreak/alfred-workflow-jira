package main

type ListFiltersCmd struct{}

func (l *ListFiltersCmd) Run(ctx *Context) error {
	ctx.wf.Run(func() {
		j, err := ctx.GetJiraClient()
		if err != nil {
			panic(err)
		}

		filters, _, err := j.Filter.GetFavouriteList()
		if err != nil {
			panic(err)
		}

		for _, filter := range filters {
			ctx.wf.NewItem(filter.Name).Arg(filter.ID).Valid(true)
		}

		ctx.wf.SendFeedback()

	})

	return nil
}
