package cn.spaceli.pgv;

import cn.spaceli.pgv.cases.TokenUse;
import org.junit.Test;

import java.util.concurrent.atomic.AtomicBoolean;

import static org.assertj.core.api.Assertions.assertThat;

public class ReflectiveValidatorIndexTest {
    @Test
    public void indexFindsOuterMessage() throws RuntimeException {
        TokenUse token = TokenUse.newBuilder().setPayload(TokenUse.Payload.newBuilder().setToken(TokenUse.Payload.Token.newBuilder().setValue("FOO"))).build();
        ReflectiveValidatorIndex index = new ReflectiveValidatorIndex();
        Validator<TokenUse> validator = index.validatorFor(TokenUse.class);

        assertThat(validator).withFailMessage("Unexpected Validator.ALWAYS_VALID").isNotEqualTo(Validator.ALWAYS_VALID);
        validator.assertValid(token);
    }

    @Test
    public void indexFindsEmbeddedMessage() throws RuntimeException {
        TokenUse.Payload payload = TokenUse.Payload.newBuilder().setToken(TokenUse.Payload.Token.newBuilder().setValue("FOO")).build();
        ReflectiveValidatorIndex index = new ReflectiveValidatorIndex();
        Validator<TokenUse.Payload> validator = index.validatorFor(TokenUse.Payload.class);

        assertThat(validator).withFailMessage("Unexpected Validator.ALWAYS_VALID").isNotEqualTo(Validator.ALWAYS_VALID);
        validator.assertValid(payload);
    }

    @Test
    public void indexFindsDoubleEmbeddedMessage() throws RuntimeException {
        TokenUse.Payload.Token token = TokenUse.Payload.Token.newBuilder().setValue("FOO").build();
        ReflectiveValidatorIndex index = new ReflectiveValidatorIndex();
        Validator<TokenUse.Payload.Token> validator = index.validatorFor(TokenUse.Payload.Token.class);

        assertThat(validator).withFailMessage("Unexpected Validator.ALWAYS_VALID").isNotEqualTo(Validator.ALWAYS_VALID);
        validator.assertValid(token);
    }

    @Test
    public void indexFallsBack() throws RuntimeException {
        AtomicBoolean called = new AtomicBoolean();
        ValidatorIndex fallback = new ValidatorIndex() {
            @Override
            @SuppressWarnings("unchecked")
            public <T> Validator<T> validatorFor(Class clazz) {
                called.set(true);
                return Validator.ALWAYS_VALID;
            }
        };

        ReflectiveValidatorIndex index = new ReflectiveValidatorIndex(fallback);
        index.validatorFor(fallback);

        assertThat(called).isTrue();
    }
}
