package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
)

// ============================================================
// 🔹 Embedding Single File
// ============================================================

//go:embed example.txt

// content:
// embedded file stored directly as string

var content string


// ============================================================
// 🔹 Embedding Directory
// ============================================================

// IMPORTANT:
//
// embed only works at:
// compile time
//
// files/folders must already exist
// BEFORE running:
//
// go run
//
// otherwise embedding fails


//go:embed NewFolder/*

// embed.FS
// → virtual read-only filesystem

var embeddedFolder embed.FS


func emdedDirective() {

	// ============================================================
	// 🔹 Reading Embedded String
	// ============================================================

	fmt.Println(
		"Embedded Content:",
		content,
	)


	// ============================================================
	// 🔹 Runtime Directory Creation
	// ============================================================

	// This directory is NOT embedded
	// because embed happens during compilation

	err := os.Mkdir(
		"RuntimeFolder",
		0755,
	)

	if err != nil {

		fmt.Println(err)

		return
	}

	defer os.RemoveAll("RuntimeFolder")


	// ============================================================
	// 🔹 Reading Embedded File
	// ============================================================

	// ReadFile()
	// → reads file from embedded filesystem

	fileContent, err := embeddedFolder.ReadFile(
		"NewFolder/text.txt",
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"\nEmbedded File Content:",
		string(fileContent),
	)


	// ============================================================
	// 🔹 Walking Embedded Filesystem
	// ============================================================

	// fs.WalkDir()
	// → recursively traverses embedded filesystem

	fmt.Println("\nWalking Embedded Filesystem:")

	err = fs.WalkDir(	embeddedFolder, "NewFolder", func(path string, d fs.DirEntry, err error ) error {

			if err != nil {

				fmt.Println(err)

				return err
			}

			fmt.Println(path)

			return nil
		})

		if err != nil {
		log.Fatal(err)
		}

// ============================================================
// 🔹 fs.WalkDir()
// ============================================================

// Syntax:
//
// fs.WalkDir(filesystem, startingPath, callbackFunction)


// Example:
//
// fs.WalkDir(folder, "NewFolder", func(...) {...})


// ------------------------------------------------------------
// 1️⃣ filesystem
// ------------------------------------------------------------

// embeddedFolder
//
// → filesystem we want to traverse
// → can be:
//    embed.FS
//    os.DirFS
//    virtual filesystem


// ------------------------------------------------------------
// 2️⃣ startingPath
// ------------------------------------------------------------

// "NewFolder"
//
// → directory/path from where traversal starts


// ------------------------------------------------------------
// 3️⃣ callbackFunction
// ------------------------------------------------------------

// func(path string, d fs.DirEntry, err error) error
//
// → automatically runs for EVERY:
//    file
//    folder
//    subfolder


// ============================================================
// 🔹 Callback Parameters
// ============================================================


// path string
//
// → current file/folder path
//
// Example:
// "NewFolder/text.txt"


// d fs.DirEntry
//
// → current directory entry
//
// Useful methods:
//
// d.Name()
// → file/folder name
//
// d.IsDir()
// → checks whether current entry is directory
//
// d.Type()
// → file type info


// err error
//
// → error during traversal
//
// if traversal fails on some file/folder,
// error comes here


// ============================================================
// 🔹 Return Value
// ============================================================

// return nil
// → continue traversal
//
// return non-nil error
// → stop traversal


// ============================================================
// 🔹 Easy Mental Model
// ============================================================

// WalkDir(
//     whereToWalk,
//     whereToStart,
//     whatToDoForEveryFile,
// )


}


// ============================================================
// 🔹 IMPORTANT EMBED CONCEPT
// ============================================================

// embed:
// → packages files into compiled binary

// Embedded files become:
// read-only assets inside executable

// Runtime-created files/directories
// are NOT embedded automatically


// ============================================================
// 🔹 QUICK REVISION
// ============================================================

 //go:embed
// → embed files/folders at compile time

// embed.FS
// → embedded virtual filesystem

// ReadFile()
// → read embedded file

// fs.WalkDir()
// → traverse embedded filesystem

// embed is:
// compile-time feature

// embedded data:
// read-only inside binary