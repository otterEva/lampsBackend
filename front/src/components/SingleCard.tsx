import { Box, Flex } from '@chakra-ui/react'

import ImageOfGood from '../UI/card_parts/ImageOfGood'
import KorzinaBtnGood from '../UI/card_parts/KorzinaBtnGood'
import PriceOfGood from '../UI/card_parts/PriceOfGood'
import TitleOfGood from '../UI/card_parts/TitleOfGood'

type SingleCardProps = {
  image_url: string
  price: string
  title: string
}

const SingleCard = ({ image_url, price, title }: SingleCardProps) => {
  return (
    <Flex
      bg="white"
      borderRadius="10px"
      flexDirection="column"
      alignItems="center"
      justifyContent="top"
      pl="10px"
      pr="10px"
    >
      <Box>
        <ImageOfGood image_url={image_url} />
      </Box>
      <Box>
        <PriceOfGood price={price} />
      </Box>
      <Box>
        <TitleOfGood title={title} />
      </Box>
      <Box>
        <KorzinaBtnGood />
      </Box>
    </Flex>
  )
}

export default SingleCard
