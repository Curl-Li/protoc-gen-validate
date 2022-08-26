package cn.spaceli.pgv;

/**
 * {@code Validator} asserts the validity of a protobuf object.
 * @param <T> The type to validate
 */
@FunctionalInterface
public interface Validator<T> {
    /**
     * Asserts validation rules on a protobuf object.
     *
     * @param proto the protobuf object to validate.
     * @throws RuntimeException with the first validation error encountered.
     */
    void assertValid(T proto) throws RuntimeException;

    /**
     * Checks validation rules on a protobuf object.
     *
     * @param proto the protobuf object to validate.
     * @return {@code true} if all rules are valid, {@code false} if not.
     */
    default boolean isValid(T proto) {
        try {
            assertValid(proto);
            return true;
        } catch (RuntimeException ex) {
            return false;
        }
    }

    Validator ALWAYS_VALID = (proto) -> {
        // Do nothing. Always valid.
    };

    Validator ALWAYS_INVALID = (proto) -> {
        throw new RuntimeException("UNKNOWN, Explicitly invalid");
    };
}
