# 排行榜组件

* 千万级数据实时排序
* 排名查询以及百分比查询
* 1000W数据15w+/s的更新操作
* 50w+/s的排名查询

## 算法

### 更新

排行榜由span数组构成，每个span由skiplists组成，链表元素为rankItem,span内元素按积分从高到低排列。span有序，因此整个span数组也是有序的。

每次更新的消耗为通过id查找item(nlog(n)) + 二分法查找span(log(n)) + span内查找插入位置的时间(常数)

### 排名查询

rankItem记录了当前所属span。根据span的下标，统计前面span的容量和 + rankItem在span中的位置(span为列表,获得位置需要链表遍历)

span的遍历累加非常耗时，因此只有前10000名(前100个span,此数值可配置)，保证完全准确的排名

对于10000以后的排名 采用span.idx x 100 + rankItem在span中的位置 获得一个大致的排名


### 百分比排名查询

rankItem数组span下标/span数组长度

### span的分裂及合并

当span大小超过预设值后，对span进行分裂，避免单个span容量太大。每隔一定的操作后，尝试将两个小的span合并为一个大的span。