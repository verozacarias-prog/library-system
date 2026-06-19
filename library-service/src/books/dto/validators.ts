import { registerDecorator, ValidationOptions } from 'class-validator';

export function IsNotFutureYear(validationOptions?: ValidationOptions) {
    return function (object: object, propertyName: string) {
        registerDecorator({
            name: 'isNotFutureYear',
            target: object.constructor,
            propertyName,
            options: validationOptions,
            validator: {
                validate(value: unknown) {
                    return typeof value === 'number' && value <= new Date().getFullYear();
                },
                defaultMessage() {
                    return `year must not be greater than ${new Date().getFullYear()}`;
                },
            },
        });
    };
}