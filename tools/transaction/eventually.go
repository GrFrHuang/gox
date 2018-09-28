// Use rabbitmq to complete transaction's eventually.
// Record place which may be general problem.
//
//BA：Basic Available 基本可用
//整个系统在某些不可抗力的情况下，仍然能够保证“可用性”，即一定时间内仍然能够返回一个明确的结果。只不过“基本可用”和“高可用”的区别是：
//“一定时间”可以适当延长
//当举行大促时，响应时间可以适当延长
//给部分用户返回一个降级页面
//给部分用户直接返回一个降级页面，从而缓解服务器压力。但要注意，返回降级页面仍然是返回明确结果。
//S：Soft State：柔性状态
//同一数据的不同副本的状态，可以不需要实时一致。
//E：Eventual Consisstency：最终一致性
//同一数据的不同副本的状态，可以不需要实时一致，但一定要保证经过一定时间后仍然是一致的。

package transaction

func BASE() {

}
