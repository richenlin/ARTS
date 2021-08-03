
/**
 * @Author: richen
 * @Date: 2020-07-15 20:53:28
 * @LastEditTime: 2020-07-15 21:54:14
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */
public class T_CAS {
    enum ReadyToRun {
        T1, T2
    }

    static volatile ReadyToRun r = ReadyToRun.T1;

    public static void main(String[] args) {
        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        new Thread(() -> {
            for (char c : aI) {
                while (r != ReadyToRun.T1) {
                }
                System.out.print(c);
                r = ReadyToRun.T2;
            }
        }, "t1").start();

        new Thread(() -> {
            for (char c : aC) {
                while (r != ReadyToRun.T2) {
                }
                System.out.print(c);
                r = ReadyToRun.T1;
            }
        }, "t2").start();
    }
}

// 自旋锁实现，自旋锁适合能够快速运算获得结果的场景