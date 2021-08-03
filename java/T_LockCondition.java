import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * @Author: richen
 * @Date: 2020-07-15 21:57:41
 * @LastEditTime: 2020-07-15 22:02:16
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */
public class T_LockCondition {
    public static void main(String[] args) {
        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        Lock lock = new ReentrantLock();

        Condition c1 = lock.newCondition();
        Condition c2 = lock.newCondition();

        new Thread(() -> {
            try {
                lock.lock();
                for (char c : aI) {
                    System.out.print(c);
                    c2.signal();
                    c1.await();
                }
                c2.signal();
            } catch (Exception e) {
                // TODO: handle exception
            } finally {
                lock.unlock();
            }

        }, "t1").start();

        new Thread(() -> {
            try {
                lock.lock();
                for (char c : aC) {
                    System.out.print(c);
                    c1.signal();
                    c2.await();
                }
                c1.signal();
            } catch (Exception e) {
                // TODO: handle exception
            } finally {
                lock.unlock();
            }
        }, "t2").start();
    }
}
// 使用可重入锁，可以精确的指定等待队列中的某一线程唤醒执行