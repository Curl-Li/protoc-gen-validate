package cn.spaceli.pgv;

import java.util.Collection;
import java.util.Map;

/**
 * {@code MapValidation} implements PGV validation for protobuf {@code Map} fields.
 */
public final class MapValidation {
    private MapValidation() {
    }

    public static void min(RuntimeException ex, Map value, int expected) {
        if (Math.min(value.size(), expected) != expected) {
            throw ex;
        }
    }

    public static void max(RuntimeException ex, Map value, int expected) {
        if (Math.max(value.size(), expected) != expected) {
            throw ex;
        }
    }

    public static void noSparse(RuntimeException ex, Map value) {
        throw new RuntimeException("no_sparse validation is not implemented for Java because protobuf maps cannot be sparse in Java");
    }

    @FunctionalInterface
    public interface MapValidator<T> {
        void accept(T val) throws RuntimeException;
    }

    public static <T> void validateParts(Collection<T> vals, MapValidator<T> validator) {
       for (T val : vals) {
           validator.accept(val);
       }
    }
}
