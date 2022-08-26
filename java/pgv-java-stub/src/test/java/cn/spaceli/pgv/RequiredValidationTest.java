package cn.spaceli.pgv;

import org.junit.Test;
import cn.spaceli.pgv.cases.Enum;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class RequiredValidationTest {
    @Test
    public void requiredWorks() throws RuntimeException {
        TestException ex = new TestException(2, "data not set");
        // Present
        RequiredValidation.required(ex, Enum.Outer.getDefaultInstance());
        // Absent
        assertThatThrownBy(() -> RequiredValidation.required(ex, null)).isEqualTo(ex);
    }
}
