package cn.spaceli.pgv;

import org.junit.Test;

import java.util.Arrays;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class RepeatedValidationTest {
    @Test
    public void minItemsWorks() throws RuntimeException {
        TestException ex = new TestException(2, "items not enough");
        // More
        RepeatedValidation.minItems(ex, Arrays.asList(10, 20, 30), 2);
        // Equal
        RepeatedValidation.minItems(ex, Arrays.asList(10, 20), 2);
        // Fewer
        assertThatThrownBy(() -> RepeatedValidation.minItems(ex, Arrays.asList(10), 2)).isNotEqualTo(Validator.ALWAYS_VALID);
    }

    @Test
    public void maxItemsWorks() throws RuntimeException {
        TestException ex = new TestException(2, "items too much");
        // More
        assertThatThrownBy(() -> RepeatedValidation.maxItems(ex, Arrays.asList(10, 20, 30), 2)).isNotEqualTo(Validator.ALWAYS_VALID);
        // Equal
        RepeatedValidation.maxItems(ex, Arrays.asList(10, 20), 2);
        // Fewer
        RepeatedValidation.maxItems(ex, Arrays.asList(10), 2);
    }

    @Test
    public void uniqueWorks() throws RuntimeException {
        TestException ex = new TestException(2, "items duplicated");
        // Unique
        RepeatedValidation.unique(ex, Arrays.asList(10, 20, 30, 40));
        // Duplicate
        assertThatThrownBy(() -> RepeatedValidation.unique(ex, Arrays.asList(10, 20, 20, 30, 30, 40))).isNotEqualTo(Validator.ALWAYS_VALID);
    }
}
