import java.util.concurrent.CountDownLatch;

/**
 * @Author: richen
 * @Date: 2020-07-15 21:21:41
 * @LastEditTime: 2020-07-15 21:54:04
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */
public class T_SyncWaitNotifys {
    public static void main(String[] args) {

        final Object o = new Object();
        CountDownLatch latcth = new CountDownLatch(1);

        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        new Thread(() -> {
            synchronized (o) {
                for (char c : aI) {
                    System.out.print(c);
                    latcth.countDown();
                    try {
                        o.notify(); // 叫醒等待队列的其他任一线程
                        o.wait(); // 让出锁
                    } catch (Exception e) {
                        // TODO: handle exception
                        e.printStackTrace();
                    }
                }
                o.notify(); // 必须，否则无法停止程序
            }
        }, "t1").start();

        new Thread(() -> {
            latcth.countDown();
            synchronized (o) {
                for (char c : aC) {
                    System.out.print(c);
                    try {
                        o.notify();
                        o.wait(); // 让出锁
                    } catch (Exception e) {
                        // TODO: handle exception
                        e.printStackTrace();
                    }
                }
                o.notify(); // 必须，否则无法停止程序
            }
        }, "t2").start();
    }
}