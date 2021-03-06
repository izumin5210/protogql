
�
google/protobuf/empty.protogoogle.protobuf"
EmptyB}
com.google.protobufB
EmptyProtoPZ.google.golang.org/protobuf/types/known/emptypb��GPB�Google.Protobuf.WellKnownTypesJ�
 3
�
 2� Protocol Buffers - Google's data interchange format
 Copyright 2008 Google Inc.  All rights reserved.
 https://developers.google.com/protocol-buffers/

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are
 met:

     * Redistributions of source code must retain the above copyright
 notice, this list of conditions and the following disclaimer.
     * Redistributions in binary form must reproduce the above
 copyright notice, this list of conditions and the following disclaimer
 in the documentation and/or other materials provided with the
 distribution.
     * Neither the name of Google Inc. nor the names of its
 contributors may be used to endorse or promote products derived from
 this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


  

" ;
	
%" ;

# E
	
# E

$ ,
	
$ ,

% +
	
% +

& "
	

& "

' !
	
$' !

( 
	
( 
�
 3 � A generic empty message that you can re-use to avoid defining duplicated
 empty messages in your APIs. A typical example is to use it as the request
 or the response type of an API method. For instance:

     service Foo {
       rpc Bar(google.protobuf.Empty) returns (google.protobuf.Empty);
     }

 The JSON representation for `Empty` is empty JSON object `{}`.



 3bproto3
�
user/user.prototestapi.usergoogle/protobuf/empty.proto"�
User
id (Rid
	full_name (	RfullName+
rolee (2.testapi.user.User.RoleRrole"'
Role
ROLE_UNSPECIFIED 	
ADMIN"N
ListUsersRequest

page_token (	R	pageToken
	page_size (RpageSize"�
ListUsersResponse(
users (2.testapi.user.UserRusers

total_size (R	totalSize&
next_page_token (	RnextPageToken"1
BatchGetUsersRequest
user_ids (RuserIds"A
BatchGetUsersResponse(
users (2.testapi.user.UserRusers"I
FollowUserRequest
user_id (RuserId
	friend_id (RfriendId"K
UnfollowUserRequest
user_id (RuserId
	friend_id (RfriendId2�
UserServiceL
	ListUsers.testapi.user.ListUsersRequest.testapi.user.ListUsersResponseE

FollowUser.testapi.user.FollowUserRequest.google.protobuf.EmptyI
UnfollowUser!.testapi.user.UnfollowUserRequest.google.protobuf.EmptyX
BatchGetUsers".testapi.user.BatchGetUsersRequest#.testapi.user.BatchGetUsersResponseBZapis/go/user;user_pbJ�

  6

  

 
	
  %

 +
	
 +


  


 

  	>

  	

  	 

  	+<

 
D

 


 
"

 
-B

 H

 

 &

 1F

 J

 

 (

 3H


  


 

  

  

  	

  

 

 

 	

 

  

  

   

   

   

  

  	

  

 

 

 

 


 




 

 

 	

 










  $


 

 !

 !


 !

 !

 !

"

"

"

"

#

#

#	

#


& (


&

 '

 '


 '

 '

 '


* ,


*

 +

 +


 +

 +

 +


. 1


.

 /

 /

 /	

 /

0

0

0	

0


3 6


3

 4

 4

 4	

 4

5

5

5	

5bproto3