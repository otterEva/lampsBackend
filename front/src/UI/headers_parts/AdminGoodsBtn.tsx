import { Button, Text } from '@chakra-ui/react'

const AdminGoodsBtn = () => {
  return (
    <Button borderRadius="10px" font-size="22px" bg="rgba(52, 72, 255, 1)" cursor="pointer">
      <Text color="white" display="inline-block" font-size="35px">
        Товары
      </Text>
    </Button>
  )
}

export default AdminGoodsBtn
