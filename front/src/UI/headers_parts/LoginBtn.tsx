import { User } from 'lucide-react'
import { useNavigate } from 'react-router-dom'

import { Button, Flex, Text } from '@chakra-ui/react'

const LoginBtn = () => {
  const navigate = useNavigate()
  return (
    <Button
      height="100%"
      borderRadius="50%"
      bg="white"
      mr="30px"
      onClick={() => navigate('/login', { replace: true })}
    >
      <Flex direction="column" align="center" justify="center">
        <User color="black" />
        <Text color="black" fontSize="sm">
          Войти
        </Text>
      </Flex>
    </Button>
  )
}

export default LoginBtn
