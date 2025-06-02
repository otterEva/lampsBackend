import { Button, Flex, Text } from '@chakra-ui/react'

type AuthBtnProps = {
  email: string
  password: string
}

const RegisterBtn = ({ email, password }: AuthBtnProps) => {
  const handleClick = async () => {
    try {
      const formData = new FormData()
      formData.append('email', email)
      formData.append('password', password)

      const response = await fetch('http://127.0.0.1:80/auth/register', {
        method: 'POST',
        body: formData,
      })

      if (response.ok) {
        alert('Регистрация прошла успешно.')
      } else {
        alert(`Ошибка`)
      }
    } catch (error) {
      console.error('Сетевая ошибка:', error)
      alert('Сетевая ошибка: не удалось отправить запрос.')
    }
  }

  return (
    <Flex width="100%" alignItems="center">
      <Button
        type="submit"
        borderRadius="10px"
        font-size="22px"
        bg="rgba(52, 72, 255, 1)"
        cursor="pointer"
        width="100%"
        onClick={handleClick}
      >
        <Text color="white" display="inline-block" font-size="35px">
          Зарегистрироваться
        </Text>
      </Button>
    </Flex>
  )
}

export default RegisterBtn
