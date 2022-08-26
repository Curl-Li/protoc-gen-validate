package cn.spaceli.pgv;

import org.junit.Test;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class ConstantValidationTest {
    @Test
    public void constantBooleanWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        ConstantValidation.constant(ex, true, true);
        assertThatThrownBy(() -> ConstantValidation.constant(ex, true, false)).isEqualTo(ex);
    }

    @Test
    public void constantFloatWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        ConstantValidation.constant(ex, 1.23F, 1.23F);
        assertThatThrownBy(() -> ConstantValidation.constant(ex, 1.23F, 3.21F)).isEqualTo(ex);
    }
}
