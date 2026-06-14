package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all courses in the Codecademy catalog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := a.effectiveLimit(0)
			courses, err := a.client.ListCourses(cmd.Context(), limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(courses, len(courses))
		},
	}
	return cmd
}

func (a *App) searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search courses by title, slug, or description",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := a.effectiveLimit(0)
			results, err := a.client.Search(cmd.Context(), args[0], limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(results, len(results))
		},
	}
	return cmd
}
