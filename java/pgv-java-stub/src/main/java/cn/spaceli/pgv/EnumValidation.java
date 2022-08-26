package cn.spaceli.pgv;

import com.google.protobuf.ProtocolMessageEnum;

/**
 * {@code EnumValidation} implements PGV validation for protobuf enumerated types.
 */
public final class EnumValidation {
    private EnumValidation() {
    }

    public static void definedOnly(RuntimeException ex, ProtocolMessageEnum value) {
        if (value.toString().equals("UNRECOGNIZED")) {
            throw ex;
        }
    }
}
