import { Flex, HStack } from '@chakra-ui/react'

import AdminGoodsBtn from '../UI/headers_parts/AdminGoodsBtn'
import AdminOrdersBtn from '../UI/headers_parts/AdminOrdersBtn'
import CatalogBtn from '../UI/headers_parts/CatalogBtn'
import KorzinaBtn from '../UI/headers_parts/KorzinaBtn'
import LoginBtn from '../UI/headers_parts/LoginBtn'
import SiriusText from '../UI/headers_parts/SiriusText'

function AdminHeader() {
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
      <Flex alignItems="center" flex="1">
        <SiriusText />
        <CatalogBtn />
      </Flex>
      <Flex alignItems="center" flex="1" justifyContent="center">
        <HStack gap="40px">
          <AdminGoodsBtn />
          <AdminOrdersBtn />
        </HStack>
      </Flex>
      <Flex alignItems="center" flex="1" justifyContent="right">
        <LoginBtn />
        <KorzinaBtn />
      </Flex>
    </Flex>
  )
}

export default AdminHeader
