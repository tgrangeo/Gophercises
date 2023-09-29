package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/spf13/cobra"
)

var listTask []string

var rootCmd = &cobra.Command{Use: "task"}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You have the following tasks:")
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("taskList"))
			b.ForEach(func(k, v []byte) error {
				fmt.Printf("%d. %s\n", k, v)
				return nil
			})
			return nil
		})
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to your TODO list",
	Run: func(cmd *cobra.Command, args []string) {
		str := ""
		for _, val := range args {
			str += val + " "
		}
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("taskList"))
			id, _ := b.NextSequence()
			err := b.Put([]byte{byte(id)}, []byte(str))
			return err
		})
		fmt.Println("Added ", str, "to your task list.")
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task from your TODO list.",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("taskList"))
			key, _ := strconv.Atoi(args[0])
			err := b.Delete([]byte{byte(key)})
			fmt.Println("TODO", key, "removed to your task list.")
			return err
		})
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all task from your TODO list.",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte("taskList"))
			return nil
		})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(resetCmd)
}

func main() {
	db, err := bolt.Open("task.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("taskList"))
		if err != nil {
			return fmt.Errorf("create bucket error: %s", err)
		}
		return nil
	})
	db.Close()
	Execute()
}
