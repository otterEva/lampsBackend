import { Flex } from '@chakra-ui/react'

import CatalogBtn from '../UI/headers_parts/CatalogBtn'
import KorzinaBtn from '../UI/headers_parts/KorzinaBtn'
import LoginBtn from '../UI/headers_parts/LoginBtn'
import SiriusText from '../UI/headers_parts/SiriusText'

function Header() {
  return (
    <Flex
      position="sticky"
      top="0"
      bg="white"
      borderBottomLeftRadius="20px"
      borderBottomRightRadius="20px"
      height="60px"
      className="header-div"
      justifyContent="space-between"
      alignItems="center"
      mb="10px"
      pl="30px"
      pr="30px"
      pt="5px"
      pb="5px"
    >
      <Flex alignItems="center" className="left-section">
        <SiriusText />
        <CatalogBtn />
      </Flex>
      <Flex alignItems="center" className="right-section">
        <LoginBtn />
        <KorzinaBtn />
      </Flex>
    </Flex>
  )
}

export default Header
