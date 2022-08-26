package cn.spaceli.pgv;

import com.google.common.primitives.Bytes;
import com.google.protobuf.ByteString;
import com.google.re2j.Pattern;

/**
 * {@code BytesValidation} implements PGV validators for protobuf {@code Byte} fields.
 */
public final class BytesValidation {
    private BytesValidation() {
    }

    public static void length(RuntimeException ex, ByteString value, int expected) {
        if (value.size() != expected) {
            throw ex;
        }
    }

    public static void minLength(RuntimeException ex, ByteString value, int expected) {
        if (value.size() < expected) {
            throw ex;
        }
    }

    public static void maxLength(RuntimeException ex, ByteString value, int expected) {
        if (value.size() > expected) {
            throw ex;
        }
    }

    public static void prefix(RuntimeException ex, ByteString value, byte[] prefix) {
        if (!value.startsWith(ByteString.copyFrom(prefix))) {
            throw ex;
        }
    }

    public static void contains(RuntimeException ex, ByteString value, byte[] contains) {
        if (Bytes.indexOf(value.toByteArray(), contains) == -1) {
            throw ex;
        }
    }

    public static void suffix(RuntimeException ex, ByteString value, byte[] suffix) {
        if (!value.endsWith(ByteString.copyFrom(suffix))) {
            throw ex;
        }
    }

    public static void pattern(RuntimeException ex, ByteString value, Pattern p) {
        if (!p.matches(value.toStringUtf8())) {
            throw ex;
        }
    }

    public static void ip(RuntimeException ex, ByteString value) {
        if (value.toByteArray().length != 4 && value.toByteArray().length != 16) {
            throw ex;
        }
    }

    public static void ipv4(RuntimeException ex, ByteString value) {
        if (value.toByteArray().length != 4) {
            throw ex;
        }
    }

    public static void ipv6(RuntimeException ex, ByteString value) {
        if (value.toByteArray().length != 16) {
            throw ex;
        }
    }
}
