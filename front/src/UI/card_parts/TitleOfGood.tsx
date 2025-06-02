import { Text } from '@chakra-ui/react'

type TitleOfGoodProps = {
  title: string
}

const TitleOfGood = ({ title }: TitleOfGoodProps) => {
  return <Text width="150px">{title}</Text>
}

export default TitleOfGood
