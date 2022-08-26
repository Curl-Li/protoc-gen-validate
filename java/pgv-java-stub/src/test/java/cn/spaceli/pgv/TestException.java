package cn.spaceli.pgv;

import org.junit.Test;

public class TestException extends RuntimeException {
    private final int code;
    private final String msg;

    public TestException(int code, String msg) {
        super(msg);
        this.code = code;
        this.msg = msg;
    }

    public int getCode() {
        return code;
    }

    public String getMsg() {
        return msg;
    }

    public static TestException UNKNOWN() {
        return new TestException(1, "UNKNOWN");
    }
}
