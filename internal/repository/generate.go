package repository

// minimock ожидает на входе пустую папку, поэтому удаляем существующую и создаем заново
//go:generate sh -c "rm -rf mocks && mkdir -p mocks"

// -i указывается интерфейс для которого генеририровать моки, ищет в текущей папке
// -o директория куда положить
// -s добавляет постфикс к имени файла
//go:generate minimock -i ChatRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i LogRepository -o ./mocks/ -s "_minimock.go"
