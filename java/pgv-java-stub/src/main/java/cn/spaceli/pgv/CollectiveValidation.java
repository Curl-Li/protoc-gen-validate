package cn.spaceli.pgv;

import java.util.Arrays;

/**
 * {@code CollectiveValidation} implements PGV validators for the collective {@code in} and {@code notIn} rules.
 */
public final class CollectiveValidation {
    private CollectiveValidation() {
    }

    public static <T> void in(RuntimeException ex, T value, T[] set) {
        for (T i : set) {
            if (value.equals(i)) {
                return;
            }
        }

        throw ex;
    }

    public static <T> void notIn(RuntimeException ex, T value, T[] set) {
        for (T i : set) {
            if (value.equals(i)) {
                throw ex;
            }
        }
    }
}
