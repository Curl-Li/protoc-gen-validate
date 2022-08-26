package cn.spaceli.pgv;

import com.google.protobuf.GeneratedMessageV3;

/**
 * {@code RequiredValidation} implements PGV validation for required fields.
 */
public final class RequiredValidation {
    private RequiredValidation() {
    }

    public static void required(RuntimeException ex, GeneratedMessageV3 value) {
        if (value == null) {
            throw ex;
        }
    }
}
