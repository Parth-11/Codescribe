package cmd

import (
	"fmt"

	"github.com/Parth-11/Codescribe/internal/ai"
	"github.com/Parth-11/Codescribe/internal/fs"

	"github.com/spf13/cobra"
)

var (
	srcDir string
	outDir string
)

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Add AI-generated comments to a codebase",
	RunE: func(cmd *cobra.Command, args []string) error {

		if srcDir == "" || outDir == "" {
			return fmt.Errorf("must provide --src and --out")
		}

		fmt.Println("📁 Copying source code...")
		files, err := fs.CopyCodebase(srcDir, outDir)
		if err != nil {
			return err
		}

		fmt.Println("🧠 Adding AI-generated comments...")

		for _, f := range files {
			if err := ai.AddCommentsToFile(f); err != nil {
				fmt.Printf("⚠️ Failed to comment %s: %v\n", f, err)
			}
		}

		fmt.Println("✅ Code commented and saved at:", outDir)
		return nil
	},
}

func init() {
	commentCmd.Flags().StringVar(&srcDir, "src", "", "source code directory")
	commentCmd.Flags().StringVar(&outDir, "out", "", "output directory for commented code")
	rootCmd.AddCommand(commentCmd)
}
