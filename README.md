# Ginkgo远程测试基础教程

[TOC]

---
## **简介**
### Go
[**Go**][1]是一款轻量级的开源语言，被设计成一门应用于搭载Web服务器，存储集群或类似用途的巨型中央服务器的系统编程语言。

### Ginkgo
Go test提供了一种简单的[自动化测试框架][2]。
[**Ginkgo**][3]以该框架为基础，构建了一个[BDD][4]测试框架，一般用于Go服务的集成测试。[**Gomega**][5]是一个匹配/断言库，Ginkgo通常与其一起使用。

## **安装**
### Go
Go可以在以下系统上运行：

- FreeBSD 10.3 or later
- Linux 2.6.23 or later with glibc
- macOS 10.8 or later
- Windows XP SP2 or later

    [*也可以在其他系统上使用源码安装*][6]

安装包下载地址为：[https://golang.org/dl/][7] 或者 [https://golang.google.cn/dl/][8]
找到对应的源码包，下载到本地： ![Installer][9]

####Linux, Mac OS X, and FreeBSD

下载后解压到`/usr/local`：

    tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
将`/usr/local/go/bin`加入环境变量中：

    export PATH=$PATH:/usr/local/go/bin

*Mac OS X也可以使用安装包直接安装到` /usr/local/go`下*

####Windows

.msi文件默认会安装在`C:\Go`目录下，将`C:\Go\bin`加入环境变量中。需要重启命令行才能生效。

####测试安装情况
创建工作目录` %USERPROFILE%\go`，创建文件`hello.go`：
``` go
package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
}
```

在当前目录下编译代码包：

    go build

生成hello.exe文件，也可以直接使用`go run hello.go`输出结果:

    hello, world

###Ginkgo和Gomega

使用`ge get`获取源码，并将可执行文件安装到`$GOPATH/bin`中：

    go get github.com/onsi/ginkgo/ginkgo
    go get github.com/onsi/gomega/...

##Ginkgo测试

首先，我们必须先有一个用以测试的对象，比如在`$GOPATH/src/books/`中新建一个文件`books.go`：

``` go
package books

type Book struct{
	Title  string;
	Author string;
	Pages  int;
}

func (book Book) CategoryByLength() string {
	if book.Pages>300{
		return "NOVEL"
	} else{
		return "SHORT STORY"
	}
}

func (book Book) GetAuthor() string  {
	return book.Author
}

```

###建立测试套件
接着，建立一个Ginkgo测试套件(test suite)，在该目录下执行命令：

    ginkgo bootstrap
    
会在当前目录下生产一个`books_suite_test.go`文件：

``` go
package books_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}
```
- 当我们执行`ginkgo`或者`go test`命令时，Go test runner会执行`TestBooks(t *testing.T)`函数。
- 当Ginkgo测试失败时，`ginkgo.Fail`会被调用，通过`gomega.RegisterFailHandler`传入Gomega中。
- 调用`RunSpecs`函数开始执行测试。

###添加测试实例

现在套件已经构建完成，但是没有测试实例，执行测试会得到空的通过信息：

    > ginkgo #or go test
    
    Running Suite: Books Suite
    ==========================
    Random Seed: 1533192357
    Will run 0 of 0 specs
    
    
    Ran 0 of 0 Specs in 0.005 seconds
    SUCCESS! -- 0 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS
    
    Ginkgo ran 1 suite in 1.511957s
    Test Suite Passed


我们向空的套件中添加测试实例，生成测试文件：
    
    ginkgo generate books

生成文件`books_test.go'：

```go
package books_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "books"
)

var _ = Describe("Books", func() {

})
```

现在向`Describe`函数中添加测试实例：

``` go
var _ = Describe("Book", func() {
	type Bookinfo struct {
		Title  string
		Author string
	}

	var (
		book 	 	Book
		bookinfo	Bookinfo
		pages 		int
		trueAuthor  string
		booklengh 	[]int
	)
	
	BeforeSuite(func ()  {
		booklengh=make([]int,2)
		booklengh[0]=1488
		booklengh[1]=24
	})
	
	BeforeEach(func() {
		bookinfo = Bookinfo{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
		}
		pages = booklengh[0]
		trueAuthor="Victor Hugo"
	})
	
	JustBeforeEach(func ()  {
		book=Book{
			Title: bookinfo.Title,
			Author: bookinfo.Author,
			Pages: pages,
		}
	})

	Describe("Categorizing book length", func () {
		Context("The first book",func ()  {
			It("should be Les Miserables",func ()  {
				Expect(book.GetAuthor()).To(Equal(trueAuthor))
			})
			It("should be a novel",func ()  {
				Expect(book.CategoryByLength()).To(Equal("NOVEL"))				
			})
		})

		Context("The second book",func ()  {
			BeforeEach(func ()  {
				bookinfo = Bookinfo{
					Title:  "Fox In Socks",
					Author: "Dr. Seuss",
				}
				pages = booklengh[1]				
			})

			It("should be Fox In Socks",func ()  {
				Expect(book.GetAuthor()).To(Equal("Dr. Seuss"))	
				Expect(book.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
	})

	Describe("Calculating total length",func ()  {
		Specify("should be more than 1000",func ()  {
			Expect(booklengh[0]+booklengh[1]).To(BeNumerically(">",1000))
		})
	})
})
```
#### Describe、Context和It
一共四个测试实例，由两个`Describe`区分，在`Describe`内部由`Context`区分。一个`It`即为一个实例。
`Describe`、`Context`和`It`内的信息即为BDD的描述，`Describe`描述个体的行为，`Context`描述其所在的情境，`It`描述所期待的结果。
`Specify`与`It`完全一样，`It`为其简写。

#### BeforeSuite
`BeforeSuite`在该套件中的测试执行前执行，通常用于初始化全局信息，读取文件等操作，如果出错则不会进行接下来的测试。`AfterSuite`在该套件中的测试执行完成后执行，通常执行一些清理工作，即使测试出错或者提前结束也会执行。
```go
	BeforeSuite(func ()  {
		booklengh=make([]int,2)
		booklengh[0]=1488
		booklengh[1]=24
	})
```
`BeforeSuite`和`AfterSuite`也可以写在`_suite_test.go`文件中，如：
```go
func TestBooks(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Books Suite")
}

var _ = AfterSuite(func() {
    dbClient.Cleanup()
    dbRunner.Stop()
})
```

#### BeforeEach
`BeforeEach`在每个测试实例开始前执行，通常用于保证每个测试实例有相同的初始状态。同一个层级内的`BeforeEach`会顺序执行，不同层级的`BeforeEach`按从外到内执行。`AfterEach`在每个测试实例结束后执行，不同层级按从内到外执行，可以用于清理，也可以进行结果的判断。
```go
	BeforeEach(func() {
		bookinfo = Bookinfo{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
		}
		pages = booklengh[0]
		trueAuthor="Victor Hugo"
	})
```
#### JustBeforeEach
`JustBeforeEach`在`BeforeEach`执行完成、测试实例开始执行之前执行。`JustBeforeEach`可以让构造和赋值分离。
```go
	JustBeforeEach(func ()  {
		book=Book{
			Title: bookinfo.Title,
			Author: bookinfo.Author,
			Pages: pages,
		}
	})
```
在案例中，不同情境下需要不同的Book实体，可以运用`BeforeEach`从外到内执行的特性，对中间变量进行不同的赋值，再根据中间变量在`JustBeforeEach`中构造Book实体。

###使用Gomega
在每个测试实例中使用`Gomega`库中的`Expect`函数进行断言：
```go
Expect(ACTUAL).To(Matcher(EXPECTED))
Expect(ACTUAL).NotTo(Matcher(EXPECTED), "Not to ...")
Expect(ACTUAL).ToNot(Matcher(EXPECTED))
```
可以使用第二个参数作为补充的输出。
Matcher包含多种匹配方式，在案例中使用了`Equal(EXPECTED)`和`BeNumerically(COMPARATOR_STRING, EXPECTED)`。除此之外，还有针对错误等其他内容的匹配，以及多个匹配的组合：
```go
Expect(err).NotTo(HaveOccurred())
Expect(ACTUAL).TO(BeTrue())
Expect(ACTUAL).TO(And(MATCHER1, MATCHER2, ...))
Expect(ACTUAL).TO(Or(MATCHER1, MATCHER2, ...))
```
Gomega也支持自定义匹配。

###完成测试
通过`ginkgo -v`输出详细的测试实例执行信息。测试通过：

    > ginkgo -v
    
    Running Suite: Books Suite
    ==========================
    Random Seed: 1533195292
    Will run 4 of 4 specs
    
    Book Categorizing book length The first book
      should be Les Miserables
      D:/documents/实习/go/src/books/books_test.go:48
    +
    ------------------------------
    Book Categorizing book length The first book
      should be a novel
      D:/documents/实习/go/src/books/books_test.go:51
    +
    ------------------------------
    Book Categorizing book length The second book
      should be Fox In Socks
      D:/documents/实习/go/src/books/books_test.go:65
    +
    ------------------------------
    Book Calculating total length
      should be more than 1000
      D:/documents/实习/go/src/books/books_test.go:73
    +
    Ran 4 of 4 Specs in 0.044 seconds
    SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS
    
    Ginkgo ran 1 suite in 1.6116903s
    Test Suite Passed


改变其中一个实例，让测试失败：
```go
			It("should be Fox In Socks",func ()  {
				Expect(book.GetAuthor()).To(Equal("Dr. Seuss"))	
				Expect(book.CategoryByLength()).NotTo(Equal("SHORT STORY"))
			})
```
执行测试：

    > ginkgo
    
    Running Suite: Books Suite
    ==========================
    Random Seed: 1533198939
    Will run 4 of 4 specs
    
    ++
    ------------------------------
    + Failure [0.000 seconds]
    Book
    D:/documents/实习/go/src/books/books_test.go:10
      Categorizing book length
      D:/documents/实习/go/src/books/books_test.go:47
        The second book
        D:/documents/实习/go/src/books/books_test.go:57
          should be Fox In Socks [It]
          D:/documents/实习/go/src/books/books_test.go:66
    
          Expected
              <string>: SHORT STORY
          not to equal
              <string>: SHORT STORY
    
          D:/documents/实习/go/src/books/books_test.go:68
    ------------------------------
    +
    
    Summarizing 1 Failure:
    
    [Fail] Book Categorizing book length The second book [It] should be Fox In Socks
    D:/documents/实习/go/src/books/books_test.go:68
    
    Ran 4 of 4 Specs in 0.034 seconds
    FAIL! -- 3 Passed | 1 Failed | 0 Pending | 0 Skipped
    --- FAIL: TestBooks (0.05s)
    FAIL
    
    Ginkgo ran 1 suite in 1.6655838s
    Test Suite Failed

ginkgo会自动输出失败的测试实例的详细信息。

####标记
可以将`P`,`X`,`F`作为标记加在`Describe`, `Context`, `It`, `Measure`前。
`P`和`X`代表该容器或者实例将不会执行。
```go
Describe("outer describe", func() {
    It("A", func() { ... })
    PIt("B", func() { ... })
})//将不会执行B。
```

`F`代表在同一层内只会执行该容器或者实例，其他的同级容器或者实例将不会执行。
```go
Describe("outer describe", func() {
    It("A", func() { ... })
    FIt("B", func() { ... })
})//将不会执行A。
```

若出现了不同层级的标记冲突，将以内层的标记为准。
```go
FDescribe("outer describe", func() {
    It("A", func() { ... })
    FIt("B", func() { ... })
})//将不会执行A。
```
`P`,`X`,`F`标记在编译期就决定了该容器或者实例是否执行。可以使用`skip`在运行期跳过测试实例：
```go
It("should do something, if it can", func() {
    if !someCondition {
        Skip("special condition wasn't met")
    }
    // assertions go here
})
```
也可以使用命令行参数`--focus=REGEXP` 和 `--skip=REGEXP`来执行、跳过符合正则表达实例。用每个实例`It`的第一个参数来进行匹配。

####测试顺序
Ginkgo默认打乱实例的执行顺序，避免测试的污染。一般来说，Ginkgo只会打乱顶层容器执行顺序，内部的实例仍然按顺序执行。可以使用命令行参数`--randomizeAllSpecs`打乱所有实例的顺序。当执行多个套件的测试时，可以使用`--randomizeSuites`打乱套件的执行顺序。
使用`--seed=SEED`来重现某次随机后的结果。

####并行测试
Ginkgo支持并行运行测试实例。它通过生成多个单独的go测试进程，并为共享队列中的每个进程提供实例来实现。使用`ginkgo -p`自动开始并行测试，也可以用`ginkgo -nodes=N`来指定测试节点数并开始测试。使用`-stream=false`来屏蔽多个节点的输出。

    > ginkgo -p -stream=false
    
    Running Suite: Books Suite
    ==========================
    Random Seed: 1533200852
    Will run 4 specs
    
    Running in parallel across 7 nodes
    
    ++++
    Ran 4 of 4 Specs in 0.010 seconds
    SUCCESS! -- 4 Passed | 0 Failed | 0 Pending | 0 Skipped
    
    
    Ginkgo ran 1 suite in 1.7602931s
    Test Suite Passed


####异步测试
Go对并发有很强大的支持。在Ginkgo中也可以在测试时使用并发。
```go
It("should post to the channel, eventually", func(done Done) {
    c := make(chan string, 0)

    go DoSomething(c)
    Expect(<-c).To(ContainSubstring("Done!"))
    close(done)
}, 0.2)
```
Done为一个通道类型的接口，使用Done为该测试设置超时时间，如0.2秒。超时则视为该测试失败。

####性能测试

`Measure`是Ginkgo的性能测试实例，与`It`使用方法相同。`Measure`使用`Benchmarker`记录测试结果，并传入一个数字来确定测试次数。
```go
Measure("should get book info",func (b Benchmarker)  {
	runtime := b.Time("runtime", func() {
		output := book.GetAuthor()
		Expect(output).To(Equal(trueAuthor))
	})
	Expect(runtime.Seconds()).To(BeNumerically("<", 0.2), "GetAuthor() shouldn't take too long.")
	b.RecordValue("book pages", (float64)(book.Pages))
},10)
```
得到结果：

    Running Suite: Books Suite
    ==========================
    Random Seed: 1533202358
    Will run 1 of 1 specs
    
    + [MEASUREMENT]
    Book
    D:/documents/实习/go/src/books/books_test.go:10
      Categorizing book
      D:/documents/实习/go/src/books/books_test.go:47
        The first book
        D:/documents/实习/go/src/books/books_test.go:48
          should get book info
          D:/documents/实习/go/src/books/books_test.go:49
    
          Ran 10 samples:
          runtime:
            Fastest Time: 0.000s
            Slowest Time: 0.000s
            Average Time: 0.000s ± 0.000s
          book pages:
            Smallest: 1488.000
             Largest: 1488.000
             Average: 1488.000 ± 0.000
    ------------------------------
    
    Ran 1 of 1 Specs in 0.028 seconds
    SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS
    
    Ginkgo ran 1 suite in 1.7204006s
    Test Suite Passed

###另一个案例
Kulshekhar Kabra的文章[Getting Started with BDD in Go Using Ginkgo][10]提供了一个较为现实的案例。他用BDD的格式列出了测试的需求：
    
    Given a shopping cart
    Given 一个购物车
      initially
      在初始状态下
        it has 0 items
        它有0种商品
        it has 0 units
        它有0个商品
        the total amount is 0.00
        总价为0.00
    
      when a new item is added
      当一个新商品被放到车中
        the shopping cart has 1 more unique item than it had earlier
        购物车比之前多了一种商品
        the shopping cart has 1 more unit than it had earlier
        购物车比之前多了一个商品
        the total amount increases by item price
        总价按新的商品的价格增加
    
      when an existing item is added
      当一个已有商品被放到车中
        the shopping cart has the same number of unique items as earlier
        购物车跟之前相比商品种类没有变化
        the shopping cart has 1 more unit than it had earlier
        购物车比之前多了一个商品
        the total amount increases by item price
        总价按新的商品的价格增加
    
      that has 0 unit of item A
      没有A类商品时
        removing item A
        拿出A类商品
          should not change the number of items
          商品种类不会改变
          should not change the number of units
          商品数量不会改变
          should not change the amount
          总价不会改变
    
      that has 1 unit of item A
      有1个A类商品时
        removing 1 unit item A
        拿出1个A类商品
          should reduce the number of items by 1
          减少1种商品
          should reduce the number of units by 1
          减少1个商品
          should reduce the amount by the item price
          总价按拿出的商品的价格减少
    
      that has 2 units of item A
      有1个A类商品时
        removing 1 unit of item A
        拿出1个A类商品
          should not reduce the number of items
          商品种类不会改变
          should reduce the number of units by 1
          减少1个商品
          should reduce the amount by the item price
          总价按拿出的商品的价格减少
    
        removing 2 units of item A
        拿出2个A类商品
          should reduce the number of items by 1
          减少1种商品
          should reduce the number of units by 2
          减少2个商品
          should reduce the amount by twice the item price
          总价按拿出的商品的价格减少两次
接着用Ginkgo框架实现了这些需求。

```go
package cart_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "."
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shopping Cart Suite")
}

var _ = Describe("Shopping cart", func() {
	itemA := Item{ID: "itemA", Name: "Item A", Price: 10.20, Qty: 0}
	itemB := Item{ID: "itemB", Name: "Item B", Price: 7.66, Qty: 0}

	Context("initially", func() {
		cart := Cart{}

		It("has 0 items", func() {
			Expect(cart.TotalUniqueItems()).Should(BeZero())
		})

		It("has 0 units", func() {
			Expect(cart.TotalUnits()).Should(BeZero())
		})

		Specify("the total amount is 0.00", func() {
			Expect(cart.TotalAmount()).Should(BeZero())
		})
	})

	Context("when a new item is added", func() {
		cart := Cart{}

		originalItemCount := cart.TotalUniqueItems()
		originalUnitCount := cart.TotalUnits()
		originalAmount := cart.TotalAmount()

		cart.AddItem(itemA)

		Context("the shopping cart", func() {
			It("has 1 more unique item than it had earlier", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount + 1))
			})

			It("has 1 more unit than it had earlier", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount + 1))
			})

			Specify("total amount increases by item price", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount + itemA.Price))
			})
		})
	})

	Context("when an existing item is added", func() {
		cart := Cart{}

		cart.AddItem(itemA)

		originalItemCount := cart.TotalUniqueItems()
		originalUnitCount := cart.TotalUnits()
		originalAmount := cart.TotalAmount()

		cart.AddItem(itemA)

		Context("the shopping cart", func() {
			It("has the same number of unique items as earlier", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount))
			})

			It("has 1 more unit than it had earlier", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount + 1))
			})

			Specify("total amount increases by item price", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount + itemA.Price))
			})
		})
	})

	Context("that has 0 unit of item A", func() {
		cart := Cart{}

		cart.AddItem(itemB) // just to mimic the existence other items
		cart.AddItem(itemB) // just to mimic the existence other items

		originalItemCount := cart.TotalUniqueItems()
		originalUnitCount := cart.TotalUnits()
		originalAmount := cart.TotalAmount()

		Context("removing item A", func() {
			cart.RemoveItem(itemA.ID, 1)

			It("should not change the number of items", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount))
			})
			It("should not change the number of units", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount))
			})
			It("should not change the amount", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount))
			})
		})
	})

	Context("that has 1 unit of item A", func() {
		cart := Cart{}

		cart.AddItem(itemB) // just to mimic the existence other items
		cart.AddItem(itemB) // just to mimic the existence other items

		cart.AddItem(itemA)

		originalItemCount := cart.TotalUniqueItems()
		originalUnitCount := cart.TotalUnits()
		originalAmount := cart.TotalAmount()

		Context("removing 1 unit item A", func() {
			cart.RemoveItem(itemA.ID, 1)

			It("should reduce the number of items by 1", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount - 1))
			})

			It("should reduce the number of units by 1", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount - 1))
			})

			It("should reduce the amount by item price", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount - itemA.Price))
			})
		})
	})

	Context("that has 2 units of item A", func() {

		Context("removing 1 unit of item A", func() {
			cart := Cart{}

			cart.AddItem(itemB) // just to mimic the existence other items
			cart.AddItem(itemB) // just to mimic the existence other items
			//Reset the cart with 2 units of item A
			cart.AddItem(itemA)
			cart.AddItem(itemA)

			originalItemCount := cart.TotalUniqueItems()
			originalUnitCount := cart.TotalUnits()
			originalAmount := cart.TotalAmount()

			cart.RemoveItem(itemA.ID, 1)

			It("should not reduce the number of items", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount))
			})

			It("should reduce the number of units by 1", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount - 1))
			})

			It("should reduce the amount by the item price", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount - itemA.Price))
			})
		})

		Context("removing 2 units of item A", func() {
			cart := Cart{}

			cart.AddItem(itemB) // just to mimic the existence other items
			cart.AddItem(itemB) // just to mimic the existence other items
			//Reset the cart with 2 units of item A
			cart.AddItem(itemA)
			cart.AddItem(itemA)

			originalItemCount := cart.TotalUniqueItems()
			originalUnitCount := cart.TotalUnits()
			originalAmount := cart.TotalAmount()

			cart.RemoveItem(itemA.ID, 2)

			It("should reduce the number of items by 1", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount - 1))
			})

			It("should reduce the number of units by 2", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount - 2))
			})

			It("should reduce the amount by twice the item price", func() {
				Expect(cart.TotalAmount()).Should(Equal(originalAmount - 2*itemA.Price))
			})
		})

	})
})
```
##使用ssh进行远程测试

[demo代码][11]

###读入文件
使用load.go：

- `Load`读取文件，返回一个`Reader`类型的变量。
- `Readln`从这个`Reader`变量中读取一行并返回改行的信息。
- `GetFile`从文件中读取每一行的信息，组成一个`string`类型的`slice`并返回。
- `GetJsonFile`从Json格式的文件中读取，并根据其初始化一组远程shell。
- 也可以使用`New`函数初始化远程shell。
###远程连接

- `connect`函数将远程shell与服务器连接，返回一个session。一个session只能接受一行命令。
###执行命令

- `Command`函数得到当前远程shell所需执行的命令集，传入`Run`函数，并行执行命令，建立多个通道接受`Run`的结果，并把返回信息打印出来。
- `start`函数调用`connect`使远程shell与服务器建立连接，并发送命令到服务器并接受返回结果，包括输出、用时、返回值等信息。
- `Run`函数为`start`提供计时，若超时则返回失败信息。

- `DoRun`函数直接调用`start`执行命令。
- 在实际开发中，可以使用New一个新的远程shell，用DoRun测试单条指令的结果；也可以使用Command读取命令集，得到返回的结果作为正确结果。可以实现读取结果文件，用作标准进行测试。

###运行结果
对两个服务器分别执行命令集

    cephmgmtclient list-soft-version
    cephmgmtclient list-soft-g
    cephmgmtclient get-license
和

    echo hello!
    ls
    435
得到结果：
```Json
{
	"10.100.47.169:22": [
		{
			"command": "echo hello!",
			"return status": "0",
			"time": "205.4499ms"
		},
		{
			"command": "ls",
			"return status": "0",
			"time": "210.4328ms"
		},
		{
			"command": "435",
			"return status": "127",
			"time": "239.3655ms"
		}
	],
	"10.100.47.76:22": [
		{
			"command": "cephmgmtclient list-soft-g",
			"return status": "2",
			"time": "2.44928s"
		},
		{
			"command": "cephmgmtclient get-license",
			"return status": "0",
			"time": "2.9025975s"
		},
		{
			"command": "cephmgmtclient list-soft-version",
			"return status": "0",
			"time": "3.0512063s"
		}
	]
}
```
###测试结果
    Running Suite: Ssh Suite
    ========================
    Random Seed: 1533205553
    Will run 3 of 3 specs
    
    Ssh Connecting to remote server and get basic info The first server
      Should get version
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:80
    
    + Failure [1.678 seconds]
    Ssh
    D:/documents/实习/go/src/gotest/ssh/ssh_test.go:13
      Connecting to remote server and get basic info
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:78
        The first server
        D:/documents/实习/go/src/gotest/ssh/ssh_test.go:79
          Should get version [It]
          D:/documents/实习/go/src/gotest/ssh/ssh_test.go:80
    
          should be less than 1.5s
          Expected
              <float64>: 1.6769284
          to be <
              <float64>: 1.5
    
          D:/documents/实习/go/src/gotest/ssh/ssh_test.go:85
    ------------------------------
    Ssh Connecting to remote server and get basic info The second server
      Should be more than 0.1s
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:97
    +
    ------------------------------
    Ssh Connecting to remote server and get basic info The second server
      Should echo twice
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:100
    +
    
    Summarizing 1 Failure:
    
    [Fail] Ssh Connecting to remote server and get basic info The first server [It] Should get version
    D:/documents/实习/go/src/gotest/ssh/ssh_test.go:85
    
    Ran 3 of 3 Specs in 2.215 seconds
    FAIL! -- 2 Passed | 1 Failed | 0 Pending | 0 Skipped
    --- FAIL: TestSsh (2.23s)
    FAIL
    
    Ginkgo ran 1 suite in 3.935117s
    Test Suite Failed
    
延长时间限制：

    Running Suite: Ssh Suite
    ========================
    Random Seed: 1533205779
    Will run 3 of 3 specs
    
    Ssh Connecting to remote server and get basic info The first server
      Should get version
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:80
    +
    ------------------------------
    Ssh Connecting to remote server and get basic info The second server
      Should be more than 0.1s
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:97
    +
    ------------------------------
    Ssh Connecting to remote server and get basic info The second server
      Should echo twice
      D:/documents/实习/go/src/gotest/ssh/ssh_test.go:100
    +
    Ran 3 of 3 Specs in 2.308 seconds
    SUCCESS! -- 3 Passed | 0 Failed | 0 Pending | 0 Skipped
    PASS
    
    Ginkgo ran 1 suite in 4.041054s
    Test Suite Passed
    


Reference

http://onsi.github.io/ginkgo/
http://onsi.github.io/gomega/
https://golang.org/doc/
https://blog.csdn.net/goodboynihaohao/article/details/79392500
https://www.cnblogs.com/Leo_wl/p/4780678.html
http://www.runoob.com/go/go-environment.html
https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo

  [1]: https://golang.org
  [2]: https://blog.csdn.net/code_segment/article/details/77507491
  [3]: http://onsi.github.io/ginkgo/
  [4]: https://www.cnblogs.com/Leo_wl/p/4780678.html
  [5]: http://onsi.github.io/gomega/
  [6]: https://golang.org/doc/install/source
  [7]: https://golang.org/dl/
  [8]: https://golang.google.cn/dl/
  [9]: http://www.runoob.com/wp-content/uploads/2015/06/golist.jpg
  [10]: https://semaphoreci.com/community/tutorials/getting-started-with-bdd-in-go-using-ginkgo
  [11]: https://github.com/Mrhelium/Ginkgo-demo