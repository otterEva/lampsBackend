import { ShoppingCart } from 'lucide-react'

import { Button, Flex, Text } from '@chakra-ui/react'

const KorzinaBtn = () => {
  return (
    <>
      <Button bg="white" border={0} borderRadius="20px">
        <Flex direction="column" align="center" justify="center">
          <ShoppingCart color="black" />
          <Text color="black" display="block">
            Корзина
          </Text>
        </Flex>
      </Button>
    </>
  )
}

export default KorzinaBtn
