package cn.spaceli.pgv;

import com.google.protobuf.Duration;
import com.google.protobuf.util.Durations;
import org.junit.Test;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class DurationValidationTest {
    @Test
    public void lessThanWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Less
        ComparativeValidation.lessThan(ex, Durations.fromSeconds(10), Durations.fromSeconds(20), Durations.comparator());
        // Equal
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, Durations.fromSeconds(10), Durations.fromSeconds(10), Durations.comparator())).isEqualTo(ex);
        // Greater
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, Durations.fromSeconds(20), Durations.fromSeconds(10), Durations.comparator())).isEqualTo(ex);
    }

    @Test
    public void lessThanOrEqualsWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Less
        ComparativeValidation.lessThanOrEqual(ex, Durations.fromSeconds(10), Durations.fromSeconds(20), Durations.comparator());
        // Equal
        ComparativeValidation.lessThanOrEqual(ex, Durations.fromSeconds(10), Durations.fromSeconds(10), Durations.comparator());
        // Greater
        assertThatThrownBy(() -> ComparativeValidation.lessThanOrEqual(ex, Durations.fromSeconds(20), Durations.fromSeconds(10), Durations.comparator())).isEqualTo(ex);
    }

    @Test
    public void greaterThanWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Less
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, Durations.fromSeconds(10), Durations.fromSeconds(20), Durations.comparator())).isEqualTo(ex);
        // Equal
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, Durations.fromSeconds(10), Durations.fromSeconds(10), Durations.comparator())).isEqualTo(ex);
        // Greater
        ComparativeValidation.greaterThan(ex, Durations.fromSeconds(20), Durations.fromSeconds(10), Durations.comparator());
    }

    @Test
    public void greaterThanOrEqualsWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Less
        assertThatThrownBy(() -> ComparativeValidation.greaterThanOrEqual(ex, Durations.fromSeconds(10), Durations.fromSeconds(20), Durations.comparator())).isEqualTo(ex);
        // Equal
        ComparativeValidation.greaterThanOrEqual(ex, Durations.fromSeconds(10), Durations.fromSeconds(10), Durations.comparator());
        // Greater
        ComparativeValidation.greaterThanOrEqual(ex, Durations.fromSeconds(20), Durations.fromSeconds(10), Durations.comparator());
    }

    @Test
    public void inWorks() throws RuntimeException {
        TestException ex = new TestException(2, "value not in set");
        Duration[] set = new Duration[]{TimestampValidation.toDuration(1, 0), TimestampValidation.toDuration(2, 0)};
        // In
        CollectiveValidation.in(ex, TimestampValidation.toDuration(1, 0), set);
        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, TimestampValidation.toDuration(3, 0), set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value in set");
        Duration[] set = new Duration[]{TimestampValidation.toDuration(1, 0), TimestampValidation.toDuration(2, 0)};
        // In
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, TimestampValidation.toDuration(1, 0), set)).isEqualTo(ex);
        // Not In
        CollectiveValidation.notIn(ex, TimestampValidation.toDuration(3, 0), set);
    }
}
