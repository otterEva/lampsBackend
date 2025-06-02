import { useNavigate } from 'react-router-dom'

import { Button, Text } from '@chakra-ui/react'

const CatalogBtn = () => {
  const navigate = useNavigate()
  return (
    <Button
      borderRadius="10px"
      font-size="22px"
      bg="rgba(52, 72, 255, 1)"
      cursor="pointer"
      onClick={() => navigate('/lamps', { replace: true })}
    >
      <Text color="white" display="inline-block" font-size="35px">
        Каталог
      </Text>
    </Button>
  )
}

export default CatalogBtn
