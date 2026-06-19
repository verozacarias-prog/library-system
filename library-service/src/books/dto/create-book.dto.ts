import { IsString, IsNotEmpty, IsInt, Min, IsISBN, ValidationOptions, registerDecorator } from 'class-validator';

function IsNotFutureYear(validationOptions?: ValidationOptions) {
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

export class CreateBookDto {

    @IsString()
    @IsNotEmpty()
    title: string;

    @IsString()
    @IsNotEmpty()
    author: string;

    @IsISBN()
    isbn: string;

    @IsInt()
    @Min(1450)
    @IsNotFutureYear()
    year: number;

    @IsString()
    @IsNotEmpty()
    genre: string;

    @IsInt()
    @Min(1)
    available_copies: number;
}
