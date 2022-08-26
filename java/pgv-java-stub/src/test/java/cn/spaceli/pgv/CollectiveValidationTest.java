package cn.spaceli.pgv;

import org.junit.Test;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class CollectiveValidationTest {
    @Test
    public void inWorks() throws RuntimeException {
        TestException ex = new TestException(2, "value not in set");
        String[] set = new String[]{"foo", "bar"};
        // In
        CollectiveValidation.in(ex, "foo", set);
        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, "baz", set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value in set");
        String[] set = new String[]{"foo", "bar"};
        // In
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, "foo", set)).isEqualTo(ex);
        // Not In
        CollectiveValidation.notIn(ex, "baz", set);
    }
}
