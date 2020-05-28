# 排行榜组件

* 500W以内数据
* 支持精确排名查询以及百分比查询
* 在支持的最大数据下15w/s的更新/精确排名查询
* 100w/s的百分比查询

##算法

###更新

排行榜由span数组构成，每个span由链表组成，链表元素为rankItem,每个span上限为100,span内元素按积分从高到低排列。
span有序，因此整个span数组也是有序的。

当rankItem更新时，使用二分查找查询rankItem所属span,通过插入法将rankItem插入到span内合适的位置，如果插入后span超过容量限制，
则span内最后一个item下降到后续span中，如果后续span容量也满，则新建一个span，将元素插入新建的span中，把新的span插入当前span以及其后续
span的中间。

每次更新的消耗为通过id查找item(nlog(n)) + 二分法查找span(log(n)) + span内查找插入位置的时间(常数)

###精确排名查询

rankItem记录了当前所属span。根据span的下标，统计前面span的容量和 + rankItem在span中的位置(span为列表,获得位置需要链表遍历)

###百分比排名查询

rankItem数组span下标/span数组长度







