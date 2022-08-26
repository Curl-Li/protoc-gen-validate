package cn.spaceli.pgv;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * {@code RepeatedValidation} implements PGV validators for collection-type validators.
 */
public final class RepeatedValidation {
    private RepeatedValidation() {
    }

    public static <T> void minItems(RuntimeException ex, List<T> values, int expected) {
        if (values.size() < expected) {
            throw ex;
        }
    }

    public static <T> void maxItems(RuntimeException ex, List<T> values, int expected) {
        if (values.size() > expected) {
            throw ex;
        }
    }

    public static <T> void unique(RuntimeException ex, List<T> values) {
        Set<T> seen = new HashSet<>();
        for (T value : values) {
            // Abort at the first sign of a duplicate
            if (!seen.add(value)) {
                throw ex;
            }
        }
    }

    @FunctionalInterface
    public interface ValidationConsumer<T> {
        void accept(T value) throws RuntimeException;
    }

    public static <T> void forEach(List<T> values, ValidationConsumer<T> consumer) {
        for (T value : values) {
            consumer.accept(value);
        }
    }
}
