import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;

/**
 * @Author: richen
 * @Date: 2020-07-15 21:10:30
 * @LastEditTime: 2020-07-15 21:13:35
 * @Description:
 * @Copyright (c) - <richenlin(at)gmail.com>
 */
public class T_BlockingQueue {
    static BlockingQueue<String> q1 = new ArrayBlockingQueue<>(1);
    static BlockingQueue<String> q2 = new ArrayBlockingQueue<>(1);

    public static void main(String[] args) {
        char[] aI = "1234567".toCharArray();
        char[] aC = "ABCDEFG".toCharArray();

        new Thread(() -> {
            for (char c : aI) {
                System.out.print(c);
                try {
                    q1.put("ok");
                    q2.take();
                } catch (Exception e) {
                    // TODO: handle exception
                    e.printStackTrace();
                }
            }
        }, "t1").start();

        new Thread(() -> {
            for (char c : aC) {
                System.out.print(c);
                try {
                    q2.put("ok");
                    q1.take();
                } catch (Exception e) {
                    // TODO: handle exception
                    e.printStackTrace();
                }
            }
        }, "t2").start();
    }
}