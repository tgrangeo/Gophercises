package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "make the task complete and remove her from your TODO list.",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			var str []byte
			b := tx.Bucket([]byte("taskList"))
			key, _ := strconv.Atoi(args[0])
			str = b.Get([]byte{byte(key)})
			time := time.Now()
			tt := time.Format("2006-01-02 15:04:05")
			comp := tx.Bucket([]byte("completeList"))
			err := comp.Put([]byte(tt), []byte(str))
			err = b.Delete([]byte{byte(key)})
			fmt.Println(string(str), "is now complete, congrats.")
			return err
		})
	},
}

var compCmd = &cobra.Command{
	Use:   "complete",
	Short: "List all of your complete tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Today you have complete all the following tasks:")
		db, err := bolt.Open("task.db", 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("completeList"))
			b.ForEach(func(k, v []byte) error {
				t, _ := time.Parse("2006-01-02 15:04:05", string(k))
				if t.Day() == time.Now().Day() && t.Month() == time.Now().Month() && time.Now().Year() == t.Year() {
					fmt.Println("    => ", string(v))
				}
				return nil
			})
			return nil
		})
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
		db.Update(func(tx *bolt.Tx) error {
			tx.DeleteBucket([]byte("completeList"))
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
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(compCmd)
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
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("completeList"))
		if err != nil {
			return fmt.Errorf("create bucket error: %s", err)
		}
		return nil
	})
	db.Close()
	Execute()
}
