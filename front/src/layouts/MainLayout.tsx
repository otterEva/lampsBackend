import { Outlet } from 'react-router-dom'

import { Box, Flex } from '@chakra-ui/react'

import AdminHeader from '../modules/AdminHeader'
import Footer from '../modules/Footer'
import Header from '../modules/Header'

type MainLayoutProps = {
  admin: boolean
}

const MainLayout = ({ admin }: MainLayoutProps) => {
  return (
    <Box bg="rgba(246, 246, 246, 1)">
      <Box ml="40px" mr="40px">
        <Flex direction="column" minHeight="95vh">
          <Box flex="1">
            {admin ? <AdminHeader /> : <Header />}
            <Box mt="20px">
              <Outlet />
            </Box>
          </Box>
        </Flex>
        <Footer />
      </Box>
    </Box>
  )
}

export default MainLayout
