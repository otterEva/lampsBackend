import { useState } from 'react'
import { Link } from 'react-router-dom'

import { FormControl } from '@chakra-ui/form-control'
import { Flex, Text, VStack } from '@chakra-ui/react'

import AuthBtn from '../UI/auth_parts/AuthBtn'
import EmailInput from '../UI/auth_parts/EmailInput'
import PasswordInput from '../UI/auth_parts/PasswordInput'

type LoginFormProps = {
  setAdmin: (admin: boolean) => void
}

function LoginInput({ setAdmin }: LoginFormProps) {
  const [data, setData] = useState({ email: '', password: '' })

  function handleFromSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const userData = {
      email: data.email,
      password: data.password,
    }
    console.log(userData)
  }

  function handleInputChange(e: React.ChangeEvent<HTMLInputElement>, cred: string) {
    setData({ ...data, [cred]: e.target.value })
  }

  return (
    <Flex alignItems="center" direction="column" mt="150px">
      <form onSubmit={handleFromSubmit}>
        <Flex
          alignItems="center"
          direction="column"
          borderRadius="25px"
          borderColor="rgba(52, 72, 255, 1)"
          borderStyle="solid"
          borderWidth="2px"
          padding="20px"
          width="300px"
          bg="white"
        >
          <FormControl width="100%">
            <VStack gap="15px">
              <Flex alignItems="center">
                <Text fontWeight="bold" fontSize="17px">
                  Введите логин и пароль
                </Text>
              </Flex>
              <EmailInput email={data.email} changeFunc={handleInputChange} />
              <PasswordInput password={data.password} changeFunc={handleInputChange} />
              <AuthBtn setAdmin={setAdmin} email={data.email} password={data.password} />
            </VStack>
          </FormControl>
          <Link to="/register">
            <Text mt="2" fontSize="11px" fontWeight="extralight">
              Ещё не зарегистрированы?
            </Text>
          </Link>
        </Flex>
      </form>
    </Flex>
  )
}

export default LoginInput
