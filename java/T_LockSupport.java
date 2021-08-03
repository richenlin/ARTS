import java.util.concurrent.CountDownLatch;
import java.util.concurrent.locks.LockSupport;

/**
 * @Author: richen
 * @Date: 2020-07-15 20:25:34
 * @LastEditTime: 2020-07-15 21:51:30
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */

public class T_LockSupport {
    static Thread t1 = null, t2 = null;

    public static void main(String[] args) {
        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        t1 = new Thread(() -> {
            for (char c : aI) {
                System.out.print(c);
                LockSupport.unpark(t2);
                LockSupport.park();
            }

        }, "t1");

        t2 = new Thread(() -> {
            for (char c : aC) {
                LockSupport.park();
                System.out.print(c);
                LockSupport.unpark(t1);
            }
        }, "t2");

        t1.start();
        t2.start();
    }
}

// 线程执行本身是无序的，t2首先被阻塞，那么不管它先执行还是后执行，都会产生一个结果就是t1先输出
// 正常来说，先阻塞线程LockSupport.park，然后被唤醒LockSupport.unpark。但是LockSupport支持先申明唤醒，后申明阻塞，那么在线程运行的时候，被先唤醒的线程是忽略阻塞的