# jewel-inject
## 描述
  1. 通过inject实现结构体间依赖管理
  2. IOC主体功能,解决单例依赖问题
  
## 说明
```
   injector := New()
   injector.Apply(&stu, &person...) //申请服务
   injector.Inject() //依赖扫描于加载
   p := injector.Service(&Person{}).(Person) //获取注册服务
```
 