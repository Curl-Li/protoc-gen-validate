package cn.spaceli.pgv;

import org.junit.Test;
import cn.spaceli.pgv.cases.Enum;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class EnumValidationTest {
    @Test
    public void definedOnlyWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Defined
        EnumValidation.definedOnly(ex, Enum.TestEnum.ONE);
        // Not Defined
        assertThatThrownBy(() -> EnumValidation.definedOnly(ex, Enum.TestEnum.UNRECOGNIZED)).isEqualTo(ex);
    }

    @Test
    public void inWorks() throws RuntimeException {
        TestException ex = new TestException(2, "value not in set");
        Enum.TestEnum[] set = new Enum.TestEnum[]{
                Enum.TestEnum.forNumber(0),
                Enum.TestEnum.forNumber(2),
        };
        // In
        CollectiveValidation.in(ex, Enum.TestEnum.TWO, set);
        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, Enum.TestEnum.ONE, set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        TestException ex = new TestException(2, "value in set");
        Enum.TestEnum[] set = new Enum.TestEnum[]{
                Enum.TestEnum.forNumber(0),
                Enum.TestEnum.forNumber(2),
        };
        // In
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, Enum.TestEnum.TWO, set)).isEqualTo(ex);
        // Not In
        CollectiveValidation.notIn(ex, Enum.TestEnum.ONE, set);
    }
}
