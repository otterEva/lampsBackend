import { jwtDecode } from 'jwt-decode'

import { Button, Flex, Text } from '@chakra-ui/react'

type AuthBtnProps = {
  email: string
  password: string
  setAdmin: (admin: boolean) => void
}

interface JwtClaims {
  userId: number
  admin: boolean
  exp: number
}

function isUserAdminFromCookie(): boolean {
  console.log('start')
  const cookies = document.cookie.split(';').map((c) => c.trim())
  console.log(cookies)
  let token = ''
  for (let c of cookies) {
    console.log(c)
    if (c.startsWith('jwt=')) {
      token = c.substring('jwt='.length)
      console.log(c)
      break
    }
  }
  if (!token) {
    console.log('no token')
    return false
  }

  try {
    const decoded = jwtDecode<JwtClaims>(token)
    console.log(decoded)
    return !!decoded.admin
  } catch (e) {
    console.error('Cannot decode JWT:', e)
    return false
  }
}

const AuthBtn = ({ email, password, setAdmin }: AuthBtnProps) => {
  const handleClick = async () => {
    try {
      const formData = new FormData()
      formData.append('email', email)
      formData.append('password', password)

      const response = await fetch('http://127.0.0.1:80/auth/login', {
        method: 'POST',
        credentials: 'include',
        body: formData,
      })

      if (response.ok) {
        alert('Авторизация прошла успешно.')
        const isadmin = isUserAdminFromCookie()
        console.log(isadmin)
        setAdmin(isadmin)
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
          Авторизоваться
        </Text>
      </Button>
    </Flex>
  )
}

export default AuthBtn
