# IOC容器以及Typescript中的实现


## IOC容器

IoC全称Inversion of Control，直译为控制反转。不是什么技术，而是一种设计思想。在OO开发中，Ioc意味着将你设计好的对象交给容器控制，而不是传统的在你的对象内部直接控制。

如何理解好Ioc呢？理解好Ioc的关键是要明确“谁控制谁，控制什么，为何是反转（有反转就应该有正转了），哪些方面反转了”，那我们来深入分析一下：

* 谁控制谁，控制什么：
  传统OO程序设计，我们直接在对象内部通过new进行创建对象，是程序主动去创建依赖对象；而IoC是有专门一个容器来创建这些对象，即由Ioc容器来控制对象的创建；谁控制谁？当然是IoC 容器控制了对象；控制什么？那就是主要控制了外部资源获取（不只是对象包括比如文件等）。

* 为何是反转，哪些方面反转了：
  有反转就有正转，传统应用程序是由我们自己在对象中主动控制去直接获取依赖对象，也就是正转；而反转则是由容器来帮忙创建及注入依赖对象；为何是反转？因为由容器帮我们查找及注入依赖对象，对象只是被动的接受依赖对象，所以是反转；哪些方面反转了？依赖对象的获取被反转了。

听着比较难以理解是不是，我们来举例说明，我们假定一个在线书店，通过BookService获取书籍：

```js
export class BookService {

  private config: DataConfig = new DataConfig();
  private dataSource: DataSource = new MysqlDataSource(config);
	
  protected constructor() {

  }

  public getBook(long bookId): Book {
      try {
          const conn = this.dataSource.getConnection();
          ...
          return book;
      } catch (err){
        throw Error("message");
      }
  }
}

```

为了从数据库查询书籍，BookService持有一个DataSource。为了实例化一个HikariDataSource，又不得不实例化一个HikariConfig。

现在，我们继续编写UserService获取用户：

```js

export class UserService {

  private config: DataConfig = new DataConfig();
  private dataSource: DataSource = new MysqlDataSource(config);

  public getUser(userId: number):User {
      try {
          const conn = this.dataSource.getConnection();
          ...
          return user;
      } catch (err){
        throw Error("message");
      }
  }
}

```
因为UserService也需要访问数据库，因此，我们不得不也实例化一个HikariDataSource。

在处理用户购买的CartController中，我们需要实例化UserService和BookService：

```js

export class CartController extends {

  private bookService = new BookService();
  private userService = new UserService(); 

  ...
}
```
类似的，在购买历史HistoryController中，也需要实例化UserService和BookService：

```js

export class HistoryController extends {

  private bookService = new BookService();
  private userService = new UserService(); 

  ...
}
```

上述每个组件都采用了一种简单的通过new创建实例并持有的方式。仔细观察，会发现以下缺点：

* 实例化一个组件，要先实例化依赖的组件，强耦合

* 每个组件都需要实例化一个依赖组件，没有复用

* 很多组件需要销毁以便释放资源，例如DataSource，但如果该组件被多个组件共享，如何确保它的使用方都已经全部被销毁

* 随着更多的组件被引入，需要共享的组件写起来会更困难，这些组件的依赖关系会越来越复杂
  

如果一个系统有大量的组件，其生命周期和相互之间的依赖关系如果由组件自身来维护，不但大大增加了系统的复杂度，而且会导致组件之间极为紧密的耦合，继而给测试和维护带来了极大的困难。

因此，核心问题是：

- 1、谁负责创建组件？
- 2、谁负责根据依赖关系组装组件？
- 3、销毁时，如何按依赖顺序正确销毁？

解决这一问题的核心方案就是IoC。

## Typescript实现IOC

参考Spring IOC的实现机制，我用Typescript实现了一个IOC容器（koatty_container），在应用启动的时候，自动分类装载组件，并且根据依赖关系，注入相应的依赖。因此，IoC又称为依赖注入（DI：Dependency Injection），它解决了一个最主要的问题：将组件的创建+配置与组件的使用相分离，并且，由IoC容器负责管理组件的生命周期。

### 组件分类

根据组件的不同应用场景，Koatty把Bean分为 'COMPONENT' | 'CONTROLLER' | 'MIDDLEWARE' | 'SERVICE' 四种类型。

* COMPONENT
  扩展类、第三方类属于此类型，例如 Plugin，ORM持久层等

* CONTROLLER
  控制器类

* MIDDLEWARE
  中间件类

* SERVICE
  逻辑服务类

### API定义

通过组件加载的Loader，在项目启动时，会自动分析并装配Bean，自动处理好Bean之间的依赖问题。IOC容器提供了一系列的API接口，方便注册以及获取装配好的Bean。

### reg<T>(target: T, options?: ObjectDefinitionOptions): T;
### reg<T>(identifier: string, target: T, options?: ObjectDefinitionOptions): T;

注册Bean到IOC容器。

* target 类或者类的实例
* identifier  别名，默认使用类名。如果自定义，从容器中获取也需要使用自定义别名
* options Bean的配置，包含作用域、生命周期、类型等等

### get(identifier: string, type?: CompomentType, args?: any[]): any;

从容器中获取Bean。

* identifier  别名，默认使用类名。如果自定义，从容器中获取也需要使用自定义别名
* type 'COMPONENT' | 'CONTROLLER' | 'MIDDLEWARE' | 'SERVICE' 四种类型。
* args 构造方法入参，如果传入参数，获取的Bean默认生命周期为Prototype，否则为单例Singleton

### getClass(identifier: string, type?: CompomentType): Function;

从容器中获取类的原型。

* identifier  别名，默认使用类名。如果自定义，从容器中获取也需要使用自定义别名
* type 'COMPONENT' | 'CONTROLLER' | 'MIDDLEWARE' | 'SERVICE' 四种类型。

### getInsByClass<T>(target: T, args?: any[]): T;

根据class类获取容器中的实例

* target 类
* args 构造方法入参，如果传入参数，获取的Bean默认生命周期为Prototype，否则为单例Singleton


----

最后，代码已经开源在[github](https://github.com/Koatty/koatty_container)，感兴趣的同学，帮忙给个star ！！