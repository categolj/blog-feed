package handler_test

import (
	. "github.com/categolj/blog-feed/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feed", func() {
	var (
		frontMatter FrontMatter
	)

	BeforeEach(func() {
		frontMatter = FrontMatter{
			Title: "Hello World",
		}
	})

	Describe("FrontMatter", func() {
		Context("Title", func() {
			It("should return hello world", func() {
				Expect(frontMatter.Title).To(Equal("Hello World"))
			})
		})
	})
})
