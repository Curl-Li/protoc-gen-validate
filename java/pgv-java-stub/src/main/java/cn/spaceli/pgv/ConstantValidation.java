package cn.spaceli.pgv;

/**
 * {@code ConstantValidation} implements PVG validators for constant values.
 */
public final class ConstantValidation {
    private ConstantValidation() {
    }

    public static <T> void constant(RuntimeException ex, T value, T expected) {
        if (!value.equals(expected)) {
            throw ex;
        }
    }
}
