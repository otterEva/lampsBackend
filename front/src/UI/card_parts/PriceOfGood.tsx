import { Text } from '@chakra-ui/react'

type PriceProps = {
  price: string
}

const Price = ({ price }: PriceProps) => {
  return (
    <Text mt="5px" mb="5px" color="rgba(70, 216, 30, 1)" fontWeight="700">
      {price} ₽
    </Text>
  )
}

export default Price
