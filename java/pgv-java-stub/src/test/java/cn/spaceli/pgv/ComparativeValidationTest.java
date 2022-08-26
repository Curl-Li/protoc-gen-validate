package cn.spaceli.pgv;

import org.junit.Test;

import java.util.Comparator;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class ComparativeValidationTest {
    @Test
    public void lessThanWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Less than
        ComparativeValidation.lessThan(ex, 10, 20, Comparator.naturalOrder());
        // Equal to
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, 10, 10, Comparator.naturalOrder())).isEqualTo(ex);
        // Greater than
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, 20, 10, Comparator.naturalOrder())).isEqualTo(ex);
    }

    @Test
    public void lessThanOrEqualWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Less than
        ComparativeValidation.lessThanOrEqual(ex, 10, 20, Comparator.naturalOrder());
        // Equal to
        ComparativeValidation.lessThanOrEqual(ex, 10, 10, Comparator.naturalOrder());
        // Greater than
        assertThatThrownBy(() -> ComparativeValidation.lessThanOrEqual(ex, 20, 10, Comparator.naturalOrder())).isEqualTo(ex);
    }

    @Test
    public void greaterThanWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Less than
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, 10, 20, Comparator.naturalOrder())).isEqualTo(ex);
        // Equal to
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, 10, 10, Comparator.naturalOrder())).isEqualTo(ex);
        // Greater than
        ComparativeValidation.greaterThan(ex, 20, 10, Comparator.naturalOrder());
    }

    @Test
    public void greaterThanOrEqualWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Less than
        assertThatThrownBy(() -> ComparativeValidation.greaterThanOrEqual(ex, 10, 20, Comparator.naturalOrder())).isEqualTo(ex);
        // Equal to
        ComparativeValidation.greaterThanOrEqual(ex, 10, 10, Comparator.naturalOrder());
        // Greater than
        ComparativeValidation.greaterThanOrEqual(ex, 20, 10, Comparator.naturalOrder());
    }

    @Test
    public void betweenInclusiveWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Lower outside
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 5, 10, true, 20, true, Comparator.naturalOrder())).isEqualTo(ex);
        // Lower bound
        ComparativeValidation.between(ex, 10, 10, true, 20, true, Comparator.naturalOrder());
        // Inside
        ComparativeValidation.between(ex, 15, 10, true, 20, true, Comparator.naturalOrder());
        // Upper bound
        ComparativeValidation.between(ex, 20, 10, true, 20, true, Comparator.naturalOrder());
        // Upper outside
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 25, 10, true, 20, true, Comparator.naturalOrder())).isEqualTo(ex);
    }

    @Test
    public void betweenExclusiveWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Lower outside
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 5, 10, false, 20, false, Comparator.naturalOrder())).isEqualTo(ex);
        // Lower bound
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 10, 10, false, 20, false, Comparator.naturalOrder())).isEqualTo(ex);
        // Inside
        ComparativeValidation.between(ex, 15, 10, false, 20, false, Comparator.naturalOrder());
        // Upper bound
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 20, 10, false, 20, false, Comparator.naturalOrder())).isEqualTo(ex);
        // Upper outside
        assertThatThrownBy(() -> ComparativeValidation.between(ex, 25, 10, false, 20, false, Comparator.naturalOrder())).isEqualTo(ex);
    }

    @Test
    public void outsideInclusiveWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Lower outside
        ComparativeValidation.outside(ex, 5, 10, true, 20, true, Comparator.naturalOrder());
        // Lower bound
        assertThatThrownBy(() -> ComparativeValidation.outside(ex, 10, 10, true, 20, true, Comparator.naturalOrder())).isEqualTo(ex);
        // Inside
        assertThatThrownBy(() -> ComparativeValidation.outside(ex, 15, 10, true, 20, true, Comparator.naturalOrder())).isEqualTo(ex);
        // Upper bound
        assertThatThrownBy(() -> ComparativeValidation.outside(ex, 20, 10, true, 20, true, Comparator.naturalOrder())).isEqualTo(ex);
        // Upper outside
        ComparativeValidation.outside(ex, 25, 10, true, 20, true, Comparator.naturalOrder());
    }

    @Test
    public void outsideExclusiveWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // Lower outside
        ComparativeValidation.outside(ex, 5, 10, false, 20, false, Comparator.naturalOrder());
        // Lower bound
        ComparativeValidation.outside(ex, 10, 10, false, 20, false, Comparator.naturalOrder());
        // Inside
        assertThatThrownBy(() -> ComparativeValidation.outside(ex, 15, 10, false, 20, false, Comparator.naturalOrder())).isEqualTo(ex);
        // Upper bound
        ComparativeValidation.outside(ex, 20, 10, false, 20, false, Comparator.naturalOrder());
        // Upper outside
        ComparativeValidation.outside(ex, 25, 10, false, 20, false, Comparator.naturalOrder());
    }

    @Test
    public void rangeChoosesCorrectly() throws RuntimeException {
        TestException ex = new TestException(1, "value not valid");
        // {gte:30, lt: 40} => x must be in the range [30, 40)
        // In between range
        ComparativeValidation.range(ex, 35, 40, null, null, 30, Comparator.naturalOrder());
        // Outside between range
        assertThatThrownBy(() -> ComparativeValidation.range(ex, 10, 40, null, null, 30, Comparator.naturalOrder())).isEqualTo(ex);

        // {lt:30, gte:40} => x must be outside the range [30, 40)
        // In outside range
        assertThatThrownBy(() -> ComparativeValidation.range(ex, 35, 30, null, null, 40, Comparator.naturalOrder())).isEqualTo(ex);
        // Outside outside range
        ComparativeValidation.range(ex, 10, 30, null, null, 40, Comparator.naturalOrder());
    }
}
