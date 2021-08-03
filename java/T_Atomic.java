
/**
 * @Author: richen
 * @Date: 2020-07-15 21:07:09
 * @LastEditTime: 2020-07-15 21:24:17
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */

import java.util.concurrent.atomic.AtomicInteger;

public class T_Atomic {
    static AtomicInteger threadNo = new AtomicInteger(1);

    public static void main(String[] args) {
        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        new Thread(() -> {
            for (char c : aI) {
                while (threadNo.get() != 1) {
                }
                System.out.print(c);
                threadNo.set(2);
            }
        }, "t1").start();

        new Thread(() -> {
            for (char c : aC) {
                while (threadNo.get() != 2) {
                }
                System.out.print(c);
                threadNo.set(1);
            }
        }, "t2").start();
    }
}

// 利用原子类型实现自旋锁